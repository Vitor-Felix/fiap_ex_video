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

func TestHandleDownload_NotFound(t *testing.T) {
	// Desliga os logs coloridos do Gin durante os testes
	gin.SetMode(gin.TestMode)

	// Setup: Instanciamos um Handler vazio (esse endpoint não usa repo nem service)
	handler := &Handler{}
	r := gin.Default()
	r.GET("/download/:filename", handler.HandleDownload)

	// Execução: Criamos uma requisição falsa para um arquivo que não existe
	req, _ := http.NewRequest(http.MethodGet, "/download/arquivo_fantasma.zip", nil)
	w := httptest.NewRecorder() // Esse cara "grava" a resposta HTTP

	r.ServeHTTP(w, req)

	// Asserção
	if w.Code != http.StatusNotFound {
		t.Errorf("Esperava status 404, recebeu %d", w.Code)
	}
}

func TestHandleDownload_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup do ambiente: Criamos um arquivo falso na pasta outputs para o teste
	utils.BasePath = "./"
	outDir := filepath.Join(utils.BasePath, "outputs")
	os.MkdirAll(outDir, 0755)
	defer os.RemoveAll(outDir) // Limpa a sujeira no final do teste!

	fakeZip := filepath.Join(outDir, "video_teste.zip")
	os.WriteFile(fakeZip, []byte("conteudo fake do zip"), 0644)

	// Setup do roteador
	handler := &Handler{}
	r := gin.Default()
	r.GET("/download/:filename", handler.HandleDownload)

	// Execução: Faz a requisição pedindo o arquivo que acabamos de criar
	req, _ := http.NewRequest(http.MethodGet, "/download/video_teste.zip", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Asserção
	if w.Code != http.StatusOK {
		t.Errorf("Esperava status 200, recebeu %d", w.Code)
	}

	// Valida se o cabeçalho de anexo foi enviado corretamente
	disposition := w.Header().Get("Content-Disposition")
	if disposition != "attachment; filename=video_teste.zip" {
		t.Errorf("Header de download incorreto, recebido: %s", disposition)
	}
}
