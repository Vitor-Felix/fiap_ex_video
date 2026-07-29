import json
import os
import time
import pika

from db import (
    update_status_to_completed,
    update_status_to_error,
    update_status_to_processing,
    get_user_email_by_video_id,
)

from email_service import send_error_email
from processor import process_video
from metrics import processed_total, start_metrics_server


RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://fiap:fiap@rabbitmq:5672/")
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
    video_id = None
    video_name = "desconhecido"

    try:
        data = json.loads(body.decode("utf-8"))

        video_id = data.get("video_id")
        video_path = data.get("video_path")
        # Extrai o nome do arquivo do path para usar nas notificações
        video_name = os.path.basename(video_path) if video_path else "desconhecido"

        print(f"📩 [MENSAGEM RECEBIDA] Video ID: {video_id} | Path: {video_path}")

        if not video_id or not video_path:
            print("❌ [WORKER] Payload inválido")
            ch.basic_nack(
                delivery_tag=method.delivery_tag,
                requeue=False,
            )
            return

        if not update_status_to_processing(video_id):
            print(f"⚠️ [RETRY] Falha ao atualizar PROCESSING {video_id}")
            ch.basic_nack(
                delivery_tag=method.delivery_tag,
                requeue=True,
            )
            return

        result = process_video(video_path, video_id)

        if result.get("success"):

            completed = update_status_to_completed(
                video_id,
                result.get("zip_path", ""),
                int(result.get("frame_count", 0)),
            )

            if completed:

                # Métrica Prometheus
                processed_total.inc()

                ch.basic_ack(
                    delivery_tag=method.delivery_tag
                )

                print(
                    f"✅ [ACK ENVIADO] Vídeo {video_id} processado com sucesso.\n"
                )

            else:
                ch.basic_nack(
                    delivery_tag=method.delivery_tag,
                    requeue=True,
                )

        else:

            error_message = str(
                result.get(
                    "error_message",
                    "Erro desconhecido",
                )
            )

            update_status_to_error(
                video_id,
                error_message,
            )

            user_email = get_user_email_by_video_id(video_id)
            send_error_email(video_id, video_name, error_message, recipient=user_email)

            ch.basic_nack(
                delivery_tag=method.delivery_tag,
                requeue=False,
            )

    except json.JSONDecodeError as e:

        print(f"❌ [WORKER] Erro JSON: {e}")

        ch.basic_nack(
            delivery_tag=method.delivery_tag,
            requeue=False,
        )

    except Exception as e:

        print(f"💥 [WORKER] Erro inesperado: {e}")

        if video_id:
            update_status_to_error(
                video_id,
                str(e),
            )
            user_email = get_user_email_by_video_id(video_id)
            send_error_email(video_id, video_name, str(e), recipient=user_email)

        ch.basic_nack(
            delivery_tag=method.delivery_tag,
            requeue=False,
        )


def start_worker():

    connection = connect_rabbitmq()

    channel = connection.channel()

    channel.queue_declare(
        queue=QUEUE_NAME,
        durable=True,
    )

    channel.basic_qos(
        prefetch_count=1
    )

    channel.basic_consume(
        queue=QUEUE_NAME,
        on_message_callback=callback,
    )

    print(
        f"🚀 [WORKER PRONTO] "
        f"Aguardando mensagens na fila '{QUEUE_NAME}'..."
    )

    try:

        channel.start_consuming()

    except KeyboardInterrupt:

        print("\n🛑 [WORKER] Encerrando consumidor...")

        channel.stop_consuming()
        connection.close()


if __name__ == "__main__":

    start_metrics_server()

    start_worker()
