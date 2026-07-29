package entities

import "time"

// User representa um usuário autenticável do sistema.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
