package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"video-processor/database"
	"video-processor/models"
	"video-processor/services"
	"video-processor/utils"
)

// HandleVideoUpload processa o upload, persiste o estado inicial no Postgres e roda o processador
func HandleVideoUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ProcessingResult{
			Success: false,
			Message: "Erro ao receber arquivo: " + err.Error(),
		})
		return
	}
	defer file.Close()

	if !services.IsValidVideoFile(header.Filename) {
		c.JSON(http.StatusBadRequest, models.ProcessingResult{
			Success: false,
			Message: "Formato de arquivo não suportado. Use: mp4, avi, mov, mkv",
		})
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s", timestamp, header.Filename)
	videoPath := filepath.Join(utils.BasePath, "uploads", filename)

	out, err := os.Create(videoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}

	// Salva a intenção de processamento no PostgreSQL como PROCESSANDO
	userID := "user_anonimo_123" // Mocado temporariamente para atender sua constraint NOT NULL
	videoID, err := database.InsertVideo(userID, header.Filename, videoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ProcessingResult{
			Success: false,
			Message: "Erro ao registrar processamento no banco: " + err.Error(),
		})
		return
	}

	// Repassa o videoID gerado (UUID) para o serviço atualizar o progresso
	result := services.ProcessVideo(videoPath, timestamp, videoID)

	if result.Success {
		os.Remove(videoPath)
	}

	c.JSON(http.StatusOK, result)
}
