package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"video-processor/dto"
	"video-processor/utils"
)

type Processor struct{}

// NewProcessor cria uma nova instância do adaptador
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessVideo executa a chamada ao FFmpeg, extrai frames e atualiza o status no repoQL via UUID
func (p *Processor) ProcessVideo(
	videoPath string,
	timestamp string,
	videoID string,
) dto.ProcessingResult {
	fmt.Printf("🎬 [ID: %s] Iniciando processamento do vídeo no FFmpeg\n", videoID)

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
		return dto.ProcessingResult{Success: false, Message: errMsg}
	}

	frames, err := filepath.Glob(filepath.Join(tempDir, "*.png"))
	if err != nil || len(frames) == 0 {
		errMsg := "Nenhum frame foi extraído do vídeo"
		return dto.ProcessingResult{Success: false, Message: errMsg}
	}

	zipFilename := fmt.Sprintf("frames_%s.zip", timestamp)
	zipPath := filepath.Join(utils.BasePath, "outputs", zipFilename)

	err = CreateZipFile(frames, zipPath)
	if err != nil {
		errMsg := "Erro ao criar arquivo ZIP: " + err.Error()
		return dto.ProcessingResult{Success: false, Message: errMsg}
	}

	imageNames := make([]string, len(frames))
	for i, frame := range frames {
		imageNames[i] = filepath.Base(frame)
	}

	return dto.ProcessingResult{
		Success:    true,
		Message:    fmt.Sprintf("Processamento concluído! %d frames extraídos.", len(frames)),
		ZipPath:    zipFilename,
		FrameCount: len(frames),
		Images:     imageNames,
	}
}
