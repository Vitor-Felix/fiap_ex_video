package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// ProcessVideo executa a chamada ao FFmpeg, extrai frames a 1fps e compacta em ZIP
func ProcessVideo(videoPath, timestamp string) models.ProcessingResult {
	fmt.Printf("🎬 Iniciando processamento: %s\n", videoPath)

	// Subindo um nível para achar a pasta "temp" na raiz do projeto
	tempDir := filepath.Join(utils.BasePath, "temp", timestamp)
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir) // Garante a limpeza da pasta temporária após o zip

	framePattern := filepath.Join(tempDir, "frame_%04d.png")

	// Dispara o subprocesso do FFmpeg nativo do S.O.
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", "fps=1",
		"-y",
		framePattern,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return models.ProcessingResult{
			Success: false,
			Message: fmt.Sprintf("Erro no ffmpeg: %s\nOutput: %s", err.Error(), string(output)),
		}
	}

	frames, err := filepath.Glob(filepath.Join(tempDir, "*.png"))
	if err != nil || len(frames) == 0 {
		return models.ProcessingResult{
			Success: false,
			Message: "Nenhum frame foi extraído do vídeo",
		}
	}

	fmt.Printf("📸 Extraídos %d frames\n", len(frames))

	zipFilename := fmt.Sprintf("frames_%s.zip", timestamp)
	zipPath := filepath.Join(utils.BasePath, "outputs", zipFilename) // Salva na pasta da raiz

	err = CreateZipFile(frames, zipPath)
	if err != nil {
		return models.ProcessingResult{
			Success: false,
			Message: "Erro ao criar arquivo ZIP: " + err.Error(),
		}
	}

	fmt.Printf("✅ ZIP criado com sucesso: %s\n", zipPath)

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
