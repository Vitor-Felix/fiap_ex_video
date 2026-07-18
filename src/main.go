package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"video-processor/database"
	"video-processor/handlers"
	"video-processor/utils"
)

func main() {
	// Garante a existência dos diretórios
	utils.CreateDirs()

	// Inicializa a conexão com o banco de dados
	database.ConnectDB()

	r := gin.Default()

	// Middleware de CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Servir arquivos estáticos (referenciando as pastas que estão um nível acima)
	r.Static("/uploads", utils.BasePath+"uploads")
	r.Static("/outputs", utils.BasePath+"outputs")

	// Rotas da aplicação delegadas para os handlers
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, handlers.GetHTMLForm())
	})

	r.POST("/upload", handlers.HandleVideoUpload)
	r.GET("/download/:filename", handlers.HandleDownload)
	r.GET("/api/status", handlers.HandleStatus)

	fmt.Println("🎬 Servidor iniciado na porta 8080")
	fmt.Println("📂 Acesse: http://localhost:8080")

	log.Fatal(r.Run(":8080"))
}
