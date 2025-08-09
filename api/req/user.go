package req

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required" validate:"min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required" validate:"min=6,max=100"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

// LogoutRequest 用户登出请求
type LogoutRequest struct {
	ID          int64  `json:"id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	LoginAt  int64  `json:"login_at,omitempty"`
	Status   int32  `json:"status,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Gender   int32  `json:"gender,omitempty"`
	Birthday string `json:"birthday,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	ID              string `json:"id" binding:"required"`
	Password        string `json:"password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// GetUserRequest 获取用户信息请求
type GetUserRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// BatchGetUsersRequest 批量获取用户信息请求
type BatchGetUsersRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
}
