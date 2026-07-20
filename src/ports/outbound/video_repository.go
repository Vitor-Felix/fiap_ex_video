package outbound

import "video-processor/domain/entities"

type VideoRepository interface {
	InsertVideo(userID, originalName, storagePath string) (string, error)

	UpdateVideoSuccess(id, zipPath string, frameCount int) error

	UpdateVideoError(id, errorMessage string) error

	GetVideosByUser(userID string) ([]entities.Video, error)
}
