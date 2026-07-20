package web

import (
	"video-processor/application"
	"video-processor/ports/outbound"
)

// Handler representa o adapter HTTP da aplicação.
type Handler struct {
	repo    outbound.VideoRepository
	service *application.VideoService
}

// NewHandler injeta as dependências necessárias.
func NewHandler(repo outbound.VideoRepository, service *application.VideoService) *Handler {
	return &Handler{
		repo:    repo,
		service: service,
	}
}
