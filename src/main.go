package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"video-processor/adapters/ffmpeg"
	"video-processor/adapters/persistence/postgres"
	"video-processor/adapters/web"
	"video-processor/application"
	"video-processor/utils"
)

func main() {
	utils.CreateDirs()

	// 1. Inicia o Adaptador de Banco de Dados
	repo, err := postgres.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Inicia o Adaptador de Processamento de Vídeo (FFmpeg)
	videoProcessor := ffmpeg.NewProcessor()

	// 3. Inicia a Regra de Negócio (Application Service), injetando as Portas (Adapters)
	videoService := application.NewVideoService(repo, videoProcessor)

	// 4. Inicia o Adaptador Web, injetando o banco e o service
	handler := web.NewHandler(repo, videoService)

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

	r.Static("/uploads", utils.BasePath+"uploads")
	r.Static("/outputs", utils.BasePath+"outputs")

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, handler.GetHTMLForm())
	})

	r.POST("/upload", handler.HandleVideoUpload)
	r.GET("/download/:filename", handler.HandleDownload)
	r.GET("/api/status", handler.HandleStatus)
	r.GET("/api/videos", handler.HandleListVideos)

	fmt.Println("🎬 Servidor iniciado na porta 8080")
	fmt.Println("📂 Acesse: http://localhost:8080")

	log.Fatal(r.Run(":8080"))
}
// CI Pipeline test
