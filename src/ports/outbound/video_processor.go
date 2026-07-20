package outbound

import "video-processor/dto"

// VideoProcessor é o contrato que qualquer motor de vídeo (FFmpeg, AWS MediaConvert, etc) deve seguir.
type VideoProcessor interface {
	ProcessVideo(videoPath string, timestamp string, videoID string) (dto.ProcessingResult, error)
}
