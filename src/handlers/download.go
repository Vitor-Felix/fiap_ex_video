package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"video-processor/utils"

	"github.com/gin-gonic/gin"
)

// HandleDownload gerencia o download dos arquivos .zip gerados
func HandleDownload(c *gin.Context) {
	filename := c.Param("filename")
	// Como o app roda dentro de src/, voltamos um nível para achar a pasta outputs na raiz
	filePath := filepath.Join(utils.BasePath, "outputs", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/zip")

	c.File(filePath)
}
