package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"video-processor/utils"

	"github.com/gin-gonic/gin"
)

func TestHandleStatus_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup: Criamos uma pasta e um zip falso para ser listado
	utils.BasePath = "./"
	outDir := filepath.Join(utils.BasePath, "outputs")
	os.MkdirAll(outDir, 0755)
	defer os.RemoveAll(outDir)

	fakeZip := filepath.Join(outDir, "teste_status.zip")
	os.WriteFile(fakeZip, []byte("conteúdo fake"), 0644)

	handler := &Handler{}
	r := gin.Default()
	r.GET("/status", handler.HandleStatus)

	// Execução: Fazemos um GET na rota
	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Asserção
	if w.Code != http.StatusOK {
		t.Errorf("Esperava status 200, recebeu %d", w.Code)
	}
}
