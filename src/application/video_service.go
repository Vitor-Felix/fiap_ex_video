package application

import (
	"fmt"
	"video-processor/dto"
	"video-processor/ports/outbound"
)

// VideoService orquestra o caso de uso de processamento de vídeos
type VideoService struct {
	repo      outbound.VideoRepository
	processor outbound.VideoProcessor
}

// NewVideoService atua como um "construtor"
func NewVideoService(repo outbound.VideoRepository, processor outbound.VideoProcessor) *VideoService {
	return &VideoService{
		repo:      repo,
		processor: processor,
	}
}

// ProcessUpload é a regra de negócio central
func (s *VideoService) ProcessUpload(userID, filename, videoPath, timestamp string) dto.ProcessingResult {
	// 1. Salva a intenção no banco (PENDENTE)
	videoID, err := s.repo.InsertVideo(userID, filename, videoPath)
	if err != nil {
		return dto.ProcessingResult{
			Success: false,
			Message: "Erro ao registrar processamento no banco: " + err.Error(),
		}
	}

	// 2. Chama o motor de processamento de forma isolada
	result, err := s.processor.ProcessVideo(videoPath, timestamp, videoID)
	if err != nil {
		// Erro de infraestrutura (ex: falha ao criar diretório)
		errorMsg := "Erro interno no processador: " + err.Error()

		if updateErr := s.repo.UpdateVideoError(videoID, errorMsg); updateErr != nil {
			fmt.Printf("⚠️ Erro ao atualizar erro no banco para ID %s: %v\n", videoID, updateErr)
		}

		return dto.ProcessingResult{
			Success: false,
			Message: errorMsg,
		}
	}

	// 3. Atualiza o banco com base no resultado do processador
	if result.Success {
		err = s.repo.UpdateVideoSuccess(videoID, result.ZipPath, result.FrameCount)
		if err != nil {
			fmt.Printf("⚠️ Erro ao atualizar sucesso no banco para ID %s: %v\n", videoID, err)
		} else {
			fmt.Printf("✅ [ID: %s] Banco atualizado com sucesso para CONCLUIDO\n", videoID)
		}
	} else {
		if updateErr := s.repo.UpdateVideoError(videoID, result.Message); updateErr != nil {
			fmt.Printf("⚠️ Erro ao atualizar erro no banco para ID %s: %v\n", videoID, updateErr)
		}
	}

	return result
}
