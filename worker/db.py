import os
import psycopg2


def get_db_connection():
    return psycopg2.connect(
        host=os.getenv("DB_HOST", "postgres"),
        port=os.getenv("DB_PORT", "5432"),
        dbname=os.getenv("DB_NAME", "fiap_x_db"),
        user=os.getenv("DB_USER", "fiap_user"),
        password=os.getenv("DB_PASSWORD", "fiap_password"),
    )


def update_status_to_processing(video_id: str) -> bool:
    """Atualiza o status do vídeo no Postgres para PROCESSANDO."""
    query = "UPDATE videos SET status = 'PROCESSANDO' WHERE id = %s;"

    conn = None
    try:
        conn = get_db_connection()
        with conn.cursor() as cursor:
            cursor.execute(query, (video_id,))
            conn.commit()
            print(f"STATUS ATUALIZADO: Vídeo {video_id} -> PROCESSANDO")
            return True
    except Exception as e:
        print(f"ERRO BANCO DE DADOS ao atualizar vídeo {video_id}: {e}")
        if conn:
            conn.rollback()
        return False
    finally:
        if conn:
            conn.close()


def update_status_to_completed(video_id: str, zip_path: str, frame_count: int) -> bool:
    query = """
        UPDATE videos
        SET status = 'CONCLUIDO',
            zip_path = %s,
            frame_count = %s,
            error_message = NULL,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = %s;
    """

    conn = None
    try:
        conn = get_db_connection()
        with conn.cursor() as cursor:
            cursor.execute(query, (zip_path, frame_count, video_id))
            conn.commit()
            print(f"STATUS ATUALIZADO: Vídeo {video_id} -> CONCLUIDO")
            return True
    except Exception as e:
        print(f"ERRO BANCO DE DADOS ao concluir vídeo {video_id}: {e}")
        if conn:
            conn.rollback()
        return False
    finally:
        if conn:
            conn.close()


def update_status_to_error(video_id: str, error_message: str) -> bool:
    query = """
        UPDATE videos
        SET status = 'ERRO',
            error_message = %s,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = %s;
    """

    conn = None
    try:
        conn = get_db_connection()
        with conn.cursor() as cursor:
            cursor.execute(query, (error_message, video_id))
            conn.commit()
            print(f"STATUS ATUALIZADO: Vídeo {video_id} -> ERRO")
            return True
    except Exception as e:
        print(f"ERRO BANCO DE DADOS ao marcar erro do vídeo {video_id}: {e}")
        if conn:
            conn.rollback()
        return False
    finally:
        if conn:
            conn.close()
