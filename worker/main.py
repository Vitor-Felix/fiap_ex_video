import json
import os
import time
import pika
from db import update_status_to_processing

RABBITMQ_URL = os.getenv(
    "RABBITMQ_URL",
    "amqp://fiap:fiap@rabbitmq:5672/"
)
QUEUE_NAME = "video_processing_queue"


def connect_rabbitmq():
    """Tenta conectar ao RabbitMQ com retry resilience para inicialização do container."""
    while True:
        try:
            print("⏳ [WORKER] Conectando ao RabbitMQ...")
            parameters = pika.URLParameters(RABBITMQ_URL)
            connection = pika.BlockingConnection(parameters)
            print("✅ [WORKER] Conectado ao RabbitMQ com sucesso!")
            return connection
        except pika.exceptions.AMQPConnectionError:
            print("⚠️ [WORKER] RabbitMQ indisponível. Retentando em 5 segundos...")
            time.sleep(5)


def callback(ch, method, properties, body):
    """Callback executado a cada nova mensagem recebida da fila."""
    try:
        data = json.loads(body.decode("utf-8"))
        video_id = data.get("video_id")
        video_path = data.get("video_path")

        print(f"📩 [MENSAGEM RECEBIDA] Video ID: {video_id} | Path: {video_path}")

        if not video_id:
            print("❌ [WORKER] Payload inválido (video_id ausente). Rejeitando mensagem.")
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
            return

        # 1. Atualiza status no banco Postgres para PROCESSANDO
        success = update_status_to_processing(video_id)

        if success:
            # 2. Envia ACK confirmando o processamento desta etapa
            ch.basic_ack(delivery_tag=method.delivery_tag)
            print(f"✅ [ACK ENVIADO] Mensagem do vídeo {video_id} processada com sucesso.\n")
        else:
            # Em caso de falha no banco, reenfileira a mensagem
            print(f"⚠️ [RETRY] Reenfileirando mensagem do vídeo {video_id}...")
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)

    except json.JSONDecodeError as e:
        print(f"❌ [WORKER] Erro ao decodificar JSON: {e}")
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
    except Exception as e:
        print(f"💥 [WORKER] Erro inesperado ao processar mensagem: {e}")
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)


def start_worker():
    connection = connect_rabbitmq()
    channel = connection.channel()

    # Garante a existência da fila
    channel.queue_declare(queue=QUEUE_NAME, durable=True)

    # Distribui apenas 1 mensagem por vez para o worker (Fair Dispatch)
    channel.basic_qos(prefetch_count=1)

    channel.basic_consume(
        queue=QUEUE_NAME,
        on_message_callback=callback
    )

    print(f"🚀 [WORKER PRONTO] Aguardando mensagens na fila '{QUEUE_NAME}'...")
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        print("\n🛑 [WORKER] Encerrando consumidor...")
        channel.stop_consuming()
        connection.close()


if __name__ == "__main__":
    start_worker()
    