package application

import (
	"fmt"
	"testing"
	"video-processor/domain/entities"
)

// ==========================================
// 1. CRIANDO OS FAKES (Mocks Manuais)
// ==========================================

// fakeRepo simula o PostgreSQL (Mantém igual)
type fakeRepo struct{}

func (f *fakeRepo) InsertVideo(userID, originalName, storagePath string) (string, error) {
	return "uuid-fake-123", nil // Simula que salvou no banco
}
func (f *fakeRepo) UpdateVideoSuccess(id, zipPath string, frameCount int) error { return nil }
func (f *fakeRepo) UpdateVideoError(id, errorMessage string) error              { return nil }
func (f *fakeRepo) GetVideosByUser(userID string) ([]entities.Video, error)     { return nil, nil }

// fakeBroker simula o RabbitMQ (Substitui o antigo fakeProcessor)
type fakeBroker struct {
	Published bool // Variável auxiliar para sabermos se o método foi chamado
}

func (f *fakeBroker) PublishVideoPending(videoID, videoPath string) error {
	f.Published = true // Marca que a mensagem foi enfileirada no teste
	return nil
}

// ==========================================
// 2. O TESTE UNITÁRIO DA REGRA DE NEGÓCIO
// ==========================================

func TestVideoService_ProcessUpload_Success(t *testing.T) {
	// Setup: Injetamos nossos fakes
	repo := &fakeRepo{}
	broker := &fakeBroker{}
	service := NewVideoService(repo, broker) // 👈 Injetamos o broker

	// Execução
	result := service.ProcessUpload("user_123", "video.mp4", "/tmp/video.mp4", "20260719")

	// Asserção (Validação)
	if !result.Success {
		t.Errorf("Esperava sucesso, mas falhou: %s", result.Message)
	}

	// Como a API não processa mais o vídeo, não checamos o ZIP.
	// Checamos se a mensagem foi pra fila!
	if !broker.Published {
		t.Errorf("Esperava que a mensagem fosse publicada no RabbitMQ, mas não foi")
	}
}

// ==========================================
// 3. TESTANDO O CAMINHO DE ERRO
// ==========================================

type fakeRepoError struct{}

func (f *fakeRepoError) InsertVideo(userID, originalName, storagePath string) (string, error) {
	return "", fmt.Errorf("banco de dados offline")
}
func (f *fakeRepoError) UpdateVideoSuccess(id, zipPath string, frameCount int) error { return nil }
func (f *fakeRepoError) UpdateVideoError(id, errorMessage string) error              { return nil }
func (f *fakeRepoError) GetVideosByUser(userID string) ([]entities.Video, error)     { return nil, nil }

func TestVideoService_ProcessUpload_Error(t *testing.T) {
	repoError := &fakeRepoError{}
	broker := &fakeBroker{}
	service := NewVideoService(repoError, broker)

	// Execução
	result := service.ProcessUpload("user_123", "video.mp4", "/tmp/video.mp4", "20260719")

	// Asserção
	if result.Success {
		t.Errorf("Esperava falha devido ao erro no banco, mas retornou sucesso")
	}

	if broker.Published {
		t.Errorf("A mensagem NÃO deveria ter ido para a fila em caso de erro no banco")
	}
}
