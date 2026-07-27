package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"video-processor/dto"
	"video-processor/utils"
)

// HandleVideoUpload processa o upload e delega a regra de negócio para o Service
func (h *Handler) HandleVideoUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ProcessingResult{
			Success: false,
			Message: "Erro ao receber arquivo: " + err.Error(),
		})
		return
	}
	// Correção 1: Função anônima com atribuição ao identificador em branco
	defer func() { _ = file.Close() }()

	// Validação inicial
	if !utils.IsValidVideoFile(header.Filename) { // 👈 Mudou de ffmpeg para utils
		c.JSON(http.StatusBadRequest, dto.ProcessingResult{
			Success: false,
			Message: "Formato de arquivo não suportado. Use: mp4, avi, mov, mkv",
		})
		return
	}

	// Salva o arquivo no disco (o Web Adapter cuida da recepção do arquivo físico)
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s", timestamp, header.Filename)
	videoPath := filepath.Join(utils.BasePath, "uploads", filename)

	out, err := os.Create(videoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}
	// Correção 2: Tratamento do out.Close()
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")

	result := h.videoService.ProcessUpload(userID, header.Filename, videoPath, timestamp)

	if result.Success {
		// O arquivo precisa continuar na pasta /uploads para o futuro Worker processá-lo!

		// 202 Accepted: A requisição foi aceita para processamento, mas ele ainda não terminou
		c.JSON(http.StatusAccepted, result)
	} else {
		// Se falhou ao salvar no banco ou no RabbitMQ, removemos o arquivo órfão
		_ = os.Remove(videoPath)
		c.JSON(http.StatusInternalServerError, result)
	}
}
