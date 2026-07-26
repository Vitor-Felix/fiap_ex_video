package dto

// LoginRequest representa o JSON recebido no endpoint de login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
