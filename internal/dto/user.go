package dto

// UserPublic 用户公开信息（不含密码）
type UserPublic struct {
	ID        int64   `json:"id" doc:"用户ID"`
	Name      string  `json:"name" doc:"用户名"`
	Email     string  `json:"email" doc:"邮箱"`
	Phone     *string `json:"phone,omitempty" doc:"手机号"`
	Status    int16   `json:"status" doc:"状态: 1-正常, 0-禁用"`
	CreatedAt string  `json:"created_at" doc:"创建时间"`
}
