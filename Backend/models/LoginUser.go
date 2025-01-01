package models

type LoginUser struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
