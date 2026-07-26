package application

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"video-processor/domain/entities"
	"video-processor/ports/outbound"
)

// AuthService contém os casos de uso relacionados à autenticação.
type AuthService struct {
	repo outbound.UserRepository
}

// NewAuthService cria um novo serviço de autenticação.
func NewAuthService(repo outbound.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// Login executa o caso de uso de autenticação.
func (s *AuthService) Login(username, password string) (*entities.User, string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, "", errors.New("usuário ou senha inválidos")
	}

	token, err := GenerateJWT(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
