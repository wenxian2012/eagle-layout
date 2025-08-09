package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-eagle/eagle-layout/api/req"
	"github.com/go-eagle/eagle-layout/api/vo"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle-layout/internal/dal/cache"
	"github.com/go-eagle/eagle-layout/internal/dal/db/model"
	"github.com/go-eagle/eagle-layout/internal/repository"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/errcode"
)

// UserHTTPService HTTP 用户服务接口
type UserHTTPService interface {
	Register(ctx context.Context, req *req.RegisterRequest) (*vo.RegisterResponse, error)
	Login(ctx context.Context, req *req.LoginRequest) (*vo.LoginResponse, error)
	Logout(ctx context.Context, req *req.LogoutRequest) (*vo.LogoutResponse, error)
	CreateUser(ctx context.Context, req *req.CreateUserRequest) (*vo.CreateUserResponse, error)
	UpdateUser(ctx context.Context, req *req.UpdateUserRequest) (*vo.UpdateUserResponse, error)
	UpdatePassword(ctx context.Context, req *req.UpdatePasswordRequest) (*vo.UpdatePasswordResponse, error)
	GetUser(ctx context.Context, req *req.GetUserRequest) (*vo.GetUserResponse, error)
	BatchGetUsers(ctx context.Context, req *req.BatchGetUsersRequest) (*vo.BatchGetUsersResponse, error)
}

type userHTTPService struct {
	repo repository.UserRepo
}

// NewUserHTTPService 创建新的 HTTP 用户服务
func NewUserHTTPService(repo repository.UserRepo) UserHTTPService {
	return &userHTTPService{
		repo: repo,
	}
}

func (s *userHTTPService) Register(ctx context.Context, req *req.RegisterRequest) (*vo.RegisterResponse, error) {
	// 检查用户是否已存在
	userBase, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}
	if userBase != nil && userBase.ID > 0 {
		return nil, errcode.ErrInvalidParam.WithDetails("User already exists")
	}

	userBase, err = s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}
	if userBase != nil && userBase.ID > 0 {
		return nil, errcode.ErrInvalidParam.WithDetails("User already exists")
	}

	// 生成密码哈希
	pwd, err := auth.HashAndSalt(req.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// 创建新用户
	user, err := s.newUser(req.Username, req.Email, pwd)
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}
	uid, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}

	// 发送欢迎邮件
	err = tasks.NewEmailWelcomeTask(tasks.EmailWelcomePayload{UserID: uid})
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}

	return &vo.RegisterResponse{
		ID:       uid,
		Username: req.Username,
	}, nil
}

func (s *userHTTPService) Login(ctx context.Context, req *req.LoginRequest) (*vo.LoginResponse, error) {
	if len(req.Email) == 0 && len(req.Username) == 0 {
		return nil, errcode.ErrInvalidParam.WithDetails("email or username is required")
	}

	// 获取用户基本信息
	var (
		user *model.UserInfoModel
		err  error
	)
	if req.Email != "" {
		user, err = s.repo.GetUserByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrInternalServer.WithDetails(err.Error())
		}
	}
	if user == nil && len(req.Username) > 0 {
		user, err = s.repo.GetUserByUsername(ctx, req.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrInternalServer.WithDetails(err.Error())
		}
	}
	if user == nil || user.ID == 0 {
		return nil, errcode.ErrUnauthorized.WithDetails("Invalid username or password")
	}

	if !auth.ComparePasswords(user.Password, req.Password) {
		return nil, errcode.ErrUnauthorized.WithDetails("Invalid username or password")
	}

	// 签名 JWT 令牌
	payload := map[string]interface{}{"user_id": user.ID, "username": user.Username}
	token, err := app.Sign(ctx, payload, app.Conf.JwtSecret, int64(cache.UserTokenExpireTime))
	if err != nil {
		return nil, errcode.ErrToken
	}

	// 将令牌记录到 Redis
	err = cache.NewUserTokenCache().SetUserTokenCache(ctx, user.ID, token, cache.UserTokenExpireTime)
	if err != nil {
		return nil, errcode.ErrToken
	}

	return &vo.LoginResponse{
		ID:          user.ID,
		AccessToken: token,
	}, nil
}

func (s *userHTTPService) Logout(ctx context.Context, req *req.LogoutRequest) (*vo.LogoutResponse, error) {
	c := cache.NewUserTokenCache()
	// 检查令牌
	token, err := c.GetUserTokenCache(ctx, req.ID)
	if err != nil {
		return nil, errcode.ErrToken
	}
	if token != req.AccessToken {
		return nil, errcode.ErrAccessDenied
	}

	// 从缓存中删除令牌
	err = c.DelUserTokenCache(ctx, req.ID)
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}

	return &vo.LogoutResponse{
		Message: "登出成功",
	}, nil
}

func (s *userHTTPService) CreateUser(ctx context.Context, req *req.CreateUserRequest) (*vo.CreateUserResponse, error) {
	// 生成密码哈希
	pwd, err := auth.HashAndSalt(req.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// 创建新用户
	user, err := s.newUser(req.Username, req.Email, pwd)
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}

	return &vo.CreateUserResponse{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
	}, nil
}

func (s *userHTTPService) UpdateUser(ctx context.Context, req *req.UpdateUserRequest) (*vo.UpdateUserResponse, error) {
	if req.UserID == 0 {
		return nil, errcode.ErrInvalidParam.WithDetails("user_id is required")
	}

	user := model.UserInfoModel{
		Nickname:  req.Nickname,
		Email:     req.Email,
		Avatar:    req.Avatar,
		Birthday:  req.Birthday,
		Bio:       req.Bio,
		Status:    req.Status,
		UpdatedAt: time.Now().Unix(),
	}
	err := s.repo.UpdateUser(ctx, req.UserID, user)
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}

	return &vo.UpdateUserResponse{
		UserID:    req.UserID,
		Nickname:  req.Nickname,
		Phone:     req.Phone,
		Email:     req.Email,
		Avatar:    req.Avatar,
		Gender:    req.Gender,
		Birthday:  req.Birthday,
		Bio:       req.Bio,
		Status:    req.Status,
		UpdatedAt: time.Now().Unix(),
	}, nil
}

func (s *userHTTPService) UpdatePassword(ctx context.Context, req *req.UpdatePasswordRequest) (*vo.UpdatePasswordResponse, error) {
	if len(req.ID) == 0 {
		return nil, errcode.ErrInvalidParam.WithDetails("id is required")
	}
	if len(req.Password) == 0 || len(req.NewPassword) == 0 || len(req.ConfirmPassword) == 0 {
		return nil, errcode.ErrInvalidParam.WithDetails("password fields are required")
	}
	if req.NewPassword != req.ConfirmPassword {
		return nil, errcode.ErrInvalidParam.WithDetails("New password and confirm password do not match")
	}

	// 获取用户基本信息
	user, err := s.repo.GetUser(ctx, cast.ToInt64(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound.WithDetails("User not found")
		}
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}
	if user == nil || user.ID == 0 {
		return nil, errcode.ErrNotFound.WithDetails("User not found")
	}

	if !auth.ComparePasswords(user.Password, req.Password) {
		return nil, errcode.ErrUnauthorized.WithDetails("Current password is incorrect")
	}

	newPwd, err := auth.HashAndSalt(req.NewPassword)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	data := model.UserInfoModel{
		Password:  newPwd,
		UpdatedAt: time.Now().Unix(),
	}
	err = s.repo.UpdateUser(ctx, user.ID, data)
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}

	return &vo.UpdatePasswordResponse{
		Message: "密码更新成功",
	}, nil
}

func (s *userHTTPService) GetUser(ctx context.Context, req *req.GetUserRequest) (*vo.GetUserResponse, error) {
	user, err := s.repo.GetUser(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound.WithDetails("User not found")
		}
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}

	u, err := s.convertUser(user)
	if err != nil {
		return nil, errcode.ErrInternalServer.WithDetails(err.Error())
	}

	return &vo.GetUserResponse{
		User: u,
	}, nil
}

func (s *userHTTPService) BatchGetUsers(ctx context.Context, req *req.BatchGetUsersRequest) (*vo.BatchGetUsersResponse, error) {
	// 检查请求是否被取消
	if ctx.Err() == context.Canceled {
		return nil, errcode.ErrDeadlineExceeded
	}

	if len(req.IDs) == 0 {
		return nil, errcode.ErrInvalidParam.WithDetails("ids is empty")
	}

	var users []*vo.User

	// 用户基本信息
	userBases, err := s.repo.BatchGetUsers(ctx, req.IDs)
	if err != nil {
		return nil, errcode.ErrInternalServer
	}
	userMap := make(map[int64]*model.UserInfoModel, 0)
	for _, val := range userBases {
		userMap[val.ID] = val
	}

	// 组织数据
	for _, id := range req.IDs {
		user, ok := userMap[id]
		if !ok {
			continue
		}
		u, err := s.convertUser(user)
		if err != nil {
			// 记录日志
			continue
		}
		users = append(users, u)
	}

	return &vo.BatchGetUsersResponse{
		Users: users,
	}, nil
}

func (s *userHTTPService) newUser(username, email, password string) (model.UserInfoModel, error) {
	return model.UserInfoModel{
		Username:  username,
		Email:     email,
		Password:  password,
		Status:    1, // 正常状态
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (s *userHTTPService) convertUser(u *model.UserInfoModel) (*vo.User, error) {
	if u == nil {
		return nil, nil
	}

	user := &vo.User{
		Id:        u.ID,
		Username:  u.Username,
		Phone:     u.Phone,
		Email:     u.Email,
		LoginAt:   u.LoginAt,
		Status:    u.Status,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Gender:    u.Gender,
		Birthday:  u.Birthday,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	return user, nil
}
