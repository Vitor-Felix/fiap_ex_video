package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"video-processor/database"
	"video-processor/models"
	"video-processor/utils"
)

// IsValidVideoFile valida se a extensão do arquivo é aceita pelo processador
func IsValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".mp4", ".avi", ".mov", ".mkv", ".wmv", ".flv", ".webm"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// ProcessVideo executa a chamada ao FFmpeg, extrai frames e atualiza o status no PostgreSQL via UUID
func ProcessVideo(videoPath, timestamp, videoID string) models.ProcessingResult {
	fmt.Printf("🎬 [ID: %s] Iniciando processamento do vídeo\n", videoID)

	tempDir := filepath.Join(utils.BasePath, "temp", timestamp)
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	framePattern := filepath.Join(tempDir, "frame_%04d.png")

	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", "fps=1",
		"-y",
		framePattern,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Erro no ffmpeg: %s\nOutput: %s", err.Error(), string(output))
		database.UpdateVideoError(videoID, errMsg) // Atualiza banco para ERRO
		return models.ProcessingResult{Success: false, Message: errMsg}
	}

	frames, err := filepath.Glob(filepath.Join(tempDir, "*.png"))
	if err != nil || len(frames) == 0 {
		errMsg := "Nenhum frame foi extraído do vídeo"
		database.UpdateVideoError(videoID, errMsg) // Atualiza banco para ERRO
		return models.ProcessingResult{Success: false, Message: errMsg}
	}

	zipFilename := fmt.Sprintf("frames_%s.zip", timestamp)
	zipPath := filepath.Join(utils.BasePath, "outputs", zipFilename)

	err = CreateZipFile(frames, zipPath)
	if err != nil {
		errMsg := "Erro ao criar arquivo ZIP: " + err.Error()
		database.UpdateVideoError(videoID, errMsg) // Atualiza banco para ERRO
		return models.ProcessingResult{Success: false, Message: errMsg}
	}

	// ✅ SE CHEGOU AQUI, DEU TUDO CERTO! Atualiza o banco para CONCLUIDO
	err = database.UpdateVideoSuccess(videoID, zipFilename, len(frames))
	if err != nil {
		fmt.Printf("⚠️ Erro ao atualizar sucesso no banco para ID %s: %v\n", videoID, err)
	}

	fmt.Printf("✅ [ID: %s] Banco atualizado com sucesso para CONCLUIDO\n", videoID)

	imageNames := make([]string, len(frames))
	for i, frame := range frames {
		imageNames[i] = filepath.Base(frame)
	}

	return models.ProcessingResult{
		Success:    true,
		Message:    fmt.Sprintf("Processamento concluído! %d frames extraídos.", len(frames)),
		ZipPath:    zipFilename,
		FrameCount: len(frames),
		Images:     imageNames,
	}
}
