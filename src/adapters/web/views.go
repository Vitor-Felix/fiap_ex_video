package web

import (
	"video-processor/utils"

	"github.com/gin-gonic/gin"
)

// ServeIndex entrega a página principal SPA a partir do sistema de arquivos estáticos
func (h *Handler) ServeIndex(c *gin.Context) {
	c.File(utils.BasePath + "web/static/index.html")
}

// RegisterStaticRoutes configura as rotas de assets estáticos (JS/CSS) no Gin
func RegisterStaticRoutes(r *gin.Engine) {
	r.Static("/static", utils.BasePath+"web/static")
}
