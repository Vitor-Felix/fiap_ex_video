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
	authService := application.NewAuthService(repo)

	// 4. Inicia o Adaptador Web, injetando os serviços
	handler := web.NewHandler(
		repo,
		videoService,
		authService,
	)

	r := gin.Default()

	// Middleware de CORS ajustado (liberando Authorization)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Servir arquivos do sistema (Uploads/Outputs)
	r.Static("/uploads", utils.BasePath+"uploads")
	r.Static("/outputs", utils.BasePath+"outputs")

	// Servir assets estáticos da interface (CSS, JS)
	web.RegisterStaticRoutes(r)

	// Rota principal servindo o HTML desacoplado
	r.GET("/", handler.ServeIndex)

	// Rota pública de Login
	r.POST("/login", handler.HandleLogin)

	// Rotas Protegidas por JWT
	r.POST(
		"/upload",
		web.AuthMiddleware(),
		handler.HandleVideoUpload,
	)

	r.GET(
		"/download/:filename",
		web.AuthMiddleware(),
		handler.HandleDownload,
	)

	r.GET(
		"/api/videos",
		web.AuthMiddleware(),
		handler.HandleListVideos,
	)

	fmt.Println("🎬 Servidor iniciado na porta 8080")
	fmt.Println("📂 Acesse: http://localhost:8080")

	log.Fatal(r.Run(":8080"))
}
