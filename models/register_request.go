package models

type RegisterRequest struct {
	Name     string `form:"name" json:"name"`
	Email    string `form:"email" json:"email"`
	Phone    string `form:"phone" json:"phone"`
	Password string `form:"password" json:"password"`
	IsAdmin  bool   `form:"is_admin" json:"is_admin"`
}