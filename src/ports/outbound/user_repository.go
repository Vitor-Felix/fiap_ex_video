package outbound

import "video-processor/domain/entities"

// UserRepository define as operações de persistência relacionadas aos usuários.
type UserRepository interface {
	GetUserByUsername(username string) (*entities.User, error)
}
