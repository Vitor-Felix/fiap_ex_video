package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"video-processor/domain/entities"

	"github.com/gin-gonic/gin"
)

// ==========================================
// 1. FAKE REPOSITORY PARA A CAMADA WEB
// ==========================================
type fakeWebRepo struct{}

func (f *fakeWebRepo) InsertVideo(u, o, s string) (string, error)  { return "", nil }
func (f *fakeWebRepo) UpdateVideoSuccess(i, z string, c int) error { return nil }
func (f *fakeWebRepo) UpdateVideoError(i, e string) error          { return nil }
func (f *fakeWebRepo) GetVideosByUser(userID string) ([]entities.Video, error) {
	// Retorna uma lista com 1 vídeo falso para o handler devolver no JSON
	return []entities.Video{
		{ID: "uuid-123", OriginalName: "video_teste.mp4", Status: "COMPLETED"},
	}, nil
}

// ==========================================
// 2. O TESTE DO HANDLER
// ==========================================
func TestHandleListVideos_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup: Instanciamos o handler injetando o repositório fake
	// Como estamos no mesmo pacote, acessamos h.repo diretamente
	handler := &Handler{
		repo: &fakeWebRepo{},
	}

	r := gin.Default()
	r.GET("/videos", handler.HandleListVideos)

	// Execução
	req, _ := http.NewRequest(http.MethodGet, "/videos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Asserção
	if w.Code != http.StatusOK {
		t.Errorf("Esperava status 200, recebeu %d", w.Code)
	}
}
