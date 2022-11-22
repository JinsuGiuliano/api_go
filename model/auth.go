package model

type Credentials struct {
	Email    string `json:"email" binding:"required" example:"user1@gmail.com"`
	Password string `json:"password" binding:"required" example:"password1"`
	// Code2FA  string `json:"code2FA"`
	// FromAPI  bool   `json:"fromApi" example:"true"`
}

type Session struct {
	Email string `json:"email"`
}
