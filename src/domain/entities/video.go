package entities

import "time"

// Video representa a entidade de domínio do processamento de vídeo.
type Video struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	OriginalName string    `json:"original_name"`
	StoragePath  string    `json:"storage_path"`
	ZipPath      string    `json:"zip_path,omitempty"`
	FrameCount   int       `json:"frame_count,omitempty"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
