package application

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"video-processor/domain/entities"
)

// fakeUserRepo implementa UserRepository para uso nos testes
type fakeUserRepo struct {
	user *entities.User
	err  error
}

func (f *fakeUserRepo) GetUserByUsername(_ string) (*entities.User, error) {
	return f.user, f.err
}

// hashPassword gera um bcrypt hash para uso nos testes
func hashPassword(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("erro ao gerar hash de senha: %v", err)
	}
	return string(h)
}

func TestAuthService_Login_Success(t *testing.T) {
	plainPassword := "senha123"
	repo := &fakeUserRepo{
		user: &entities.User{
			ID:           "uuid-abc",
			Username:     "joao",
			PasswordHash: hashPassword(t, plainPassword),
		},
	}

	svc := NewAuthService(repo)
	user, token, err := svc.Login("joao", plainPassword)

	if err != nil {
		t.Fatalf("esperava sucesso no login, mas obteve erro: %v", err)
	}
	if user == nil || user.ID != "uuid-abc" {
		t.Errorf("usuário retornado incorreto: %+v", user)
	}
	if token == "" {
		t.Error("esperava um token JWT não vazio")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	repo := &fakeUserRepo{
		user: &entities.User{
			ID:           "uuid-abc",
			Username:     "joao",
			PasswordHash: hashPassword(t, "senha_correta"),
		},
	}

	svc := NewAuthService(repo)
	user, token, err := svc.Login("joao", "senha_errada")

	if err == nil {
		t.Fatal("esperava erro de senha inválida, mas login foi bem-sucedido")
	}
	if user != nil || token != "" {
		t.Error("não deveria retornar usuário ou token em caso de senha errada")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := &fakeUserRepo{
		user: nil,
		err:  errors.New("usuário não encontrado"),
	}

	svc := NewAuthService(repo)
	user, token, err := svc.Login("nao_existe", "qualquer")

	if err == nil {
		t.Fatal("esperava erro de usuário não encontrado")
	}
	if user != nil || token != "" {
		t.Error("não deveria retornar usuário ou token quando o repositório falha")
	}
}
