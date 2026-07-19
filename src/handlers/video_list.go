package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"video-processor/database"
)

// HandleListVideos retorna o histórico de tarefas em formato JSON
func HandleListVideos(c *gin.Context) {
	userID := "user_anonimo_123" // O mesmo ID mocado do upload

	videos, err := database.GetVideosByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar histórico: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, videos)
}
