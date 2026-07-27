package messaging

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQAdapter implementa a porta MessageBroker
type RabbitMQAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// Mensagem que será trafegada na fila (Payload)
type VideoMessage struct {
	VideoID   string `json:"video_id"`
	VideoPath string `json:"video_path"`
}

// NewRabbitMQAdapter conecta no broker e declara a fila
func NewRabbitMQAdapter(url string) (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// QueueDeclare garante que a fila existe antes de mandarmos algo para ela
	q, err := ch.QueueDeclare(
		"video_processing_queue", // nome da fila
		true,                     // durable (sobrevive a restart do RabbitMQ)
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // argumentos extras
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQAdapter{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

// PublishVideoPending serializa os dados e envia para a fila
func (r *RabbitMQAdapter) PublishVideoPending(videoID, videoPath string) error {
	// Timeout de 5s para evitar que a goroutine da requisição trave se a rede cair
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := VideoMessage{
		VideoID:   videoID,
		VideoPath: videoPath,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(ctx,
		"",           // exchange (padrão)
		r.queue.Name, // routing key (nome da fila)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Garante que a mensagem seja salva no disco do RabbitMQ
			Body:         body,
		})
}
