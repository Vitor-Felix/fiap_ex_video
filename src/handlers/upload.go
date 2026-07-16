package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	// Ajuste o caminho abaixo conforme o nome do seu módulo no go.mod
	"video-processor/models"
	"video-processor/services"
	"video-processor/utils"
)

// HandleVideoUpload processa o upload e dispara a extração de frames
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

	// Salvando na pasta raiz (voltando um nível da pasta src/)
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

	// Delega o processamento pesado para a camada de serviço dedicada
	result := services.ProcessVideo(videoPath, timestamp)

	if result.Success {
		os.Remove(videoPath)
	}

	c.JSON(http.StatusOK, result)
}
