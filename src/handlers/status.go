package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"video-processor/utils"

	"github.com/gin-gonic/gin"
)

// HandleStatus lista os arquivos processados e metadados
func HandleStatus(c *gin.Context) {
	// Buscando na pasta raiz voltando um nível
	files, err := filepath.Glob(filepath.Join(utils.BasePath, "outputs", "*.zip"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar arquivos"})
		return
	}

	var results []map[string]interface{}
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"filename":     filepath.Base(file),
			"size":         info.Size(),
			"created_at":   info.ModTime().Format("2006-01-02 15:04:05"),
			"download_url": "/download/" + filepath.Base(file),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": results,
		"total": len(results),
	})
}
