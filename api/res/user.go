package res

import "time"

// User include user base info and user profile
type User struct {
	Id        int64      `json:"id"`
	Username  string     `json:"username"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Status    int32      `json:"status"`
	Nickname  string     `json:"nickname"`
	Avatar    string     `json:"avatar"`
	Gender    int32      `json:"gender"`
	Birthday  string     `json:"birthday"`
	Bio       string     `json:"bio"`
	LoginAt   *time.Time `json:"login_at"` // login time for last times
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	ID           int64  `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// LogoutResponse 用户登出响应
type LogoutResponse struct {
	Message string `json:"message"`
}

// CreateUserResponse 创建用户响应
type CreateUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UpdateUserResponse 更新用户响应
type UpdateUserResponse struct {
	UserID    int64      `json:"user_id"`
	Username  string     `json:"username,omitempty"`
	Email     string     `json:"email,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Status    int32      `json:"status,omitempty"`
	Nickname  string     `json:"nickname,omitempty"`
	Avatar    string     `json:"avatar,omitempty"`
	Gender    int32      `json:"gender,omitempty"`
	Birthday  string     `json:"birthday,omitempty"`
	Bio       string     `json:"bio,omitempty"`
	LoginAt   *time.Time `json:"login_at"` // login time for last times
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// UpdatePasswordResponse 更新密码响应
type UpdatePasswordResponse struct {
	Message string `json:"message"`
}

// GetUserResponse 获取用户信息响应
type GetUserResponse struct {
	User *User `json:"user"`
}

// BatchGetUsersResponse 批量获取用户信息响应
type BatchGetUsersResponse struct {
	Users []*User `json:"users"`
}
