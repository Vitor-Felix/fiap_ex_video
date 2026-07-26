package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleListVideos retorna o histórico de tarefas em formato JSON
func (h *Handler) HandleListVideos(c *gin.Context) {
	userID := c.GetString("user_id")

	videos, err := h.videoRepo.GetVideosByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar histórico: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, videos)
}
