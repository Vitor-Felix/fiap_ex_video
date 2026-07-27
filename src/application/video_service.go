package application

import (
	"fmt"
	"video-processor/dto"
	"video-processor/ports/outbound"
)

type VideoService struct {
	repo   outbound.VideoRepository
	broker outbound.MessageBroker
}

func NewVideoService(repo outbound.VideoRepository, broker outbound.MessageBroker) *VideoService {
	return &VideoService{
		repo:   repo,
		broker: broker,
	}
}

func (s *VideoService) ProcessUpload(userID, filename, videoPath, timestamp string) dto.ProcessingResult {
	// 1. Salva a intenção no banco (PENDENTE)
	videoID, err := s.repo.InsertVideo(userID, filename, videoPath)
	if err != nil {
		return dto.ProcessingResult{
			Success: false,
			Message: "Erro ao registrar processamento no banco: " + err.Error(),
		}
	}

	// 2. Publica o evento na mensageria
	err = s.broker.PublishVideoPending(videoID, videoPath)
	if err != nil {
		// Compensação: Se o RabbitMQ falhar, marcamos o status como ERRO no banco
		_ = s.repo.UpdateVideoError(videoID, "Falha ao enfileirar no RabbitMQ")
		return dto.ProcessingResult{
			Success: false,
			Message: "Erro ao enviar para fila de processamento: " + err.Error(),
		}
	}

	fmt.Printf("✅ [ID: %s] Vídeo PENDENTE e enviado para a fila do RabbitMQ\n", videoID)

	// 3. Retorna sucesso IMEDIATO
	return dto.ProcessingResult{
		Success: true,
		Message: "Upload concluído. O vídeo entrou na fila de processamento.",
	}
}
