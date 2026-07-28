package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"video-processor/adapters/messaging"
	"video-processor/adapters/persistence/postgres"
	"video-processor/adapters/web"
	"video-processor/application"
	"video-processor/utils"
)

func main() {
	utils.CreateDirs()

	// Registra as métricas personalizadas da aplicação
	utils.RegisterMetrics()

	// 1. Inicia o Banco de Dados
	repo, err := postgres.NewRepository()
	if err != nil {
		log.Fatal("Erro no BD:", err)
	}

	// 2. Inicia o RabbitMQ
	// Pega a URL do docker-compose ou usa localhost localmente
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://fiap:fiap@localhost:5672/"
	}

	broker, err := messaging.NewRabbitMQAdapter(rabbitURL)
	if err != nil {
		log.Fatal("Erro ao conectar no RabbitMQ: ", err)
	}

	// 3. Inicia a Regra de Negócio injetando o RabbitMQ
	videoService := application.NewVideoService(repo, broker)
	authService := application.NewAuthService(repo)

	// 4. Inicia o Adaptador Web
	handler := web.NewHandler(repo, videoService, authService)

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

	// =====================================================
	// Endpoint de métricas para o Prometheus
	// =====================================================
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
	fmt.Println("📊 Métricas Prometheus: http://localhost:8080/metrics")

	log.Fatal(r.Run(":8080"))
}
