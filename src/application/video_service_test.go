package application

import (
	"fmt"
	"testing"
	"video-processor/domain/entities"
	"video-processor/dto"
)

// ==========================================
// 1. CRIANDO OS FAKES (Mocks Manuais)
// ==========================================

// fakeRepo simula o PostgreSQL
type fakeRepo struct{}

func (f *fakeRepo) InsertVideo(userID, originalName, storagePath string) (string, error) {
	return "uuid-fake-123", nil // Simula que salvou no banco
}
func (f *fakeRepo) UpdateVideoSuccess(id, zipPath string, frameCount int) error {
	return nil
}
func (f *fakeRepo) UpdateVideoError(id, errorMessage string) error {
	return nil
}
func (f *fakeRepo) GetVideosByUser(userID string) ([]entities.Video, error) {
	return nil, nil
}

// fakeProcessor simula o FFmpeg
type fakeProcessor struct{}

func (f *fakeProcessor) ProcessVideo(videoPath, timestamp, videoID string) dto.ProcessingResult {
	// Retorna um sucesso fixo para enganar o Service
	return dto.ProcessingResult{
		Success:    true,
		Message:    "Mock concluído",
		ZipPath:    "frames_fake.zip",
		FrameCount: 10,
	}
}

// ==========================================
// 2. O TESTE UNITÁRIO DA REGRA DE NEGÓCIO
// ==========================================

func TestVideoService_ProcessUpload_Success(t *testing.T) {
	// Setup: Injetamos nossos fakes no lugar dos adaptadores reais
	repo := &fakeRepo{}
	processor := &fakeProcessor{}
	service := NewVideoService(repo, processor)

	// Execução: Chamamos a regra de negócio
	result := service.ProcessUpload("user_123", "video.mp4", "/tmp/video.mp4", "20260719")

	// Asserção (Validação)
	if !result.Success {
		t.Errorf("Esperava sucesso, mas falhou: %s", result.Message)
	}

	if result.ZipPath != "frames_fake.zip" {
		t.Errorf("Esperava o zip 'frames_fake.zip', mas recebeu '%s'", result.ZipPath)
	}
}

// ==========================================
// 3. TESTANDO O CAMINHO DE ERRO
// ==========================================

// fakeRepoError simula um banco de dados que está fora do ar
type fakeRepoError struct{}

func (f *fakeRepoError) InsertVideo(userID, originalName, storagePath string) (string, error) {
	return "", fmt.Errorf("banco de dados offline") // 👈 Forçando o erro aqui
}
func (f *fakeRepoError) UpdateVideoSuccess(id, zipPath string, frameCount int) error { return nil }
func (f *fakeRepoError) UpdateVideoError(id, errorMessage string) error              { return nil }
func (f *fakeRepoError) GetVideosByUser(userID string) ([]entities.Video, error)     { return nil, nil }

func TestVideoService_ProcessUpload_Error(t *testing.T) {
	// Setup injetando o repo que falha
	repoError := &fakeRepoError{}
	processor := &fakeProcessor{}
	service := NewVideoService(repoError, processor)

	// Execução
	result := service.ProcessUpload("user_123", "video.mp4", "/tmp/video.mp4", "20260719")

	// Asserção
	if result.Success {
		t.Errorf("Esperava falha devido ao erro no banco, mas retornou sucesso")
	}
}
