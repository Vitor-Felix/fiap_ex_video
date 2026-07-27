import os
import psycopg2


def get_db_connection():
    return psycopg2.connect(
        host=os.getenv("DB_HOST", "postgres"),
        port=os.getenv("DB_PORT", "5432"),
        dbname=os.getenv("DB_NAME", "fiap_db"),
        user=os.getenv("DB_USER", "postgres"),
        password=os.getenv("DB_PASSWORD", "postgres")
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
            