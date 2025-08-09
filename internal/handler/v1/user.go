package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle-layout/api/req"
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/errcode"
)

// UserHandler 用户
type UserHandler struct {
	UserHTTPService service.UserHTTPService
}

// NewUserHandler 创建新的用户
func NewUserHandler(userHTTPService service.UserHTTPService) *UserHandler {
	return &UserHandler{
		UserHTTPService: userHTTPService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户使用用户名、邮箱和密码进行注册
// @Tags user
// @Accept json
// @Produce json
// @Param request body req.RegisterRequest true "注册请求"
// @Success 200 {object} req.RegisterResponse "注册成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "内部服务器错误"
// @Router /v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req req.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	resp, err := h.UserHTTPService.Register(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, err)
		return
	}

	app.Success(c, resp)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户使用用户名/邮箱和密码进行登录
// @Tags user
// @Accept json
// @Produce json
// @Param request body req.LoginRequest true "登录请求"
// @Success 200 {object} req.LoginResponse "登录成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 401 {object} error "认证失败"
// @Failure 500 {object} error "内部服务器错误"
// @Router /v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req req.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	resp, err := h.UserHTTPService.Login(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, err)
		return
	}

	app.Success(c, resp)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户退出登录
// @Tags user
// @Accept json
// @Produce json
// @Param request body req.LogoutRequest true "登出请求"
// @Success 200 {object} req.LogoutResponse "登出成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 401 {object} error "认证失败"
// @Failure 500 {object} error "内部服务器错误"
// @Router /v1/auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	var req req.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	resp, err := h.UserHTTPService.Logout(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, err)
		return
	}

	app.Success(c, resp)
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 管理员创建新用户
// @Tags user
// @Accept json
// @Produce json
// @Param request body req.CreateUserRequest true "创建用户请求"
// @Success 200 {object} req.CreateUserResponse "创建成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "内部服务器错误"
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req req.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	resp, err := h.UserHTTPService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, err)
		return
	}

	app.Success(c, resp)
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} req.GetUserResponse "获取成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 404 {object} error "用户不存在"
// @Failure 500 {object} error "内部服务器错误"
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails("invalid user id"))
		return
	}

	req := &req.GetUserRequest{ID: id}
	resp, err := h.UserHTTPService.GetUser(c.Request.Context(), req)
	if err != nil {
		app.Error(c, err)
		return
	}

	app.Success(c, resp)
}

// BatchGetUsers 批量获取用户信息
// @Summary 批量获取用户信息
// @Description 根据用户ID列表批量获取用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param request body req.BatchGetUsersRequest true "批量获取用户请求"
// @Success 200 {object} req.BatchGetUsersResponse "获取成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "内部服务器错误"
// @Router /v1/users/batch [post]
func (h *UserHandler) BatchGetUsers(c *gin.Context) {
	var req req.BatchGetUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	resp, err := h.UserHTTPService.BatchGetUsers(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, err)
		return
	}

	app.Success(c, resp)
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户的基本信息
// @Tags user
// @Accept json
// @Produce json
// @Param request body req.UpdateUserRequest true "更新用户请求"
// @Success 200 {object} req.UpdateUserResponse "更新成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 404 {object} error "用户不存在"
// @Failure 500 {object} error "内部服务器错误"
// @Router /v1/users [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req req.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	resp, err := h.UserHTTPService.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, err)
		return
	}

	app.Success(c, resp)
}

// UpdatePassword 更新用户密码
// @Summary 更新用户密码
// @Description 用户更新自己的密码
// @Tags user
// @Accept json
// @Produce json
// @Param request body req.UpdatePasswordRequest true "更新密码请求"
// @Success 200 {object} req.UpdatePasswordResponse "更新成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 401 {object} error "原密码错误"
// @Failure 404 {object} error "用户不存在"
// @Failure 500 {object} error "内部服务器错误"
// @Router /v1/users/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req req.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	resp, err := h.UserHTTPService.UpdatePassword(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, err)
		return
	}

	app.Success(c, resp)
}
