package models

import "time"

// Video representa a entidade exata da sua tabela com UUID e ENUM
type Video struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	OriginalName string    `json:"original_name"`
	StoragePath  string    `json:"storage_path"`
	ZipPath      string    `json:"zip_path,omitempty"`
	FrameCount   int       `json:"frame_count,omitempty"`
	Status       string    `json:"status"` // PENDENTE, PROCESSANDO, CONCLUIDO, ERRO
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type VideoRequest struct {
	VideoPath string `json:"video_path"`
	OutputDir string `json:"output_dir"`
}

type ProcessingResult struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	ZipPath    string   `json:"zip_path,omitempty"`
	FrameCount int      `json:"frame_count,omitempty"`
	Images     []string `json:"images,omitempty"`
}
