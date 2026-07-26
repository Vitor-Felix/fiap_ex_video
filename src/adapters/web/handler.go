package web

import (
	"video-processor/application"
	"video-processor/ports/outbound"
)

// Handler representa o adapter HTTP da aplicação.
type Handler struct {
	videoRepo    outbound.VideoRepository
	videoService *application.VideoService
	authService  *application.AuthService
}

// NewHandler injeta as dependências necessárias.
func NewHandler(
	videoRepo outbound.VideoRepository,
	videoService *application.VideoService,
	authService *application.AuthService,
) *Handler {
	return &Handler{
		videoRepo:    videoRepo,
		videoService: videoService,
		authService:  authService,
	}
}
