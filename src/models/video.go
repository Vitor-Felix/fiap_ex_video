package models

// VideoRequest representa a estrutura para requisições diretas de processamento
type VideoRequest struct {
	VideoPath string `json:"video_path"`
	OutputDir string `json:"output_dir"`
}

// ProcessingResult representa a resposta padrão após o processamento do vídeo
type ProcessingResult struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	ZipPath    string   `json:"zip_path,omitempty"`
	FrameCount int      `json:"frame_count,omitempty"`
	Images     []string `json:"images,omitempty"`
}
