package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"video-processor/dto"
)

// HandleLogin autentica um usuário.
// Nesta primeira versão apenas verifica se o usuário existe.
func (h *Handler) HandleLogin(c *gin.Context) {
	var request dto.LoginRequest

	// Converte o JSON recebido para a struct LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON inválido: " + err.Error(),
		})
		return
	}

	user, token, err := h.authService.Login(
		request.Username,
		request.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário ou senha inválidos.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login realizado com sucesso.",
		"username": user.Username,
		"token":    token,
	})
}
