package outbound

// MessageBroker define o contrato para envio de mensagens assíncronas
type MessageBroker interface {
	PublishVideoPending(videoID, videoPath string) error
}
