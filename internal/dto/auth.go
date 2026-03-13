package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	Name     string `json:"name" required:"true" minLength:"1" maxLength:"100" doc:"用户名"`
	Email    string `json:"email" required:"true" format:"email" doc:"邮箱"`
	Phone    string `json:"phone" maxLength:"20" doc:"手机号"`
	Password string `json:"password" required:"true" minLength:"6" doc:"密码"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" required:"true" format:"email" doc:"邮箱"`
	Password string `json:"password" required:"true" doc:"密码"`
}

// AuthResponse 认证响应（登录/注册）
type AuthResponse struct {
	Token string     `json:"token" doc:"JWT Token"`
	User  UserPublic `json:"user" doc:"用户信息"`
}
