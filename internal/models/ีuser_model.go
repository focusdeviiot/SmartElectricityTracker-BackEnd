package models

type GetUserRes struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     Role   `json:"role"`
}
