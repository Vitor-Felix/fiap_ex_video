package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"video-processor/models"

	_ "github.com/lib/pq" // Driver do PostgreSQL anotado de forma anônima
)

// DB é a instância global de conexão com o banco de dados
var DB *sql.DB

// ConnectDB inicializa a conexão com o PostgreSQL usando variáveis de ambiente
func ConnectDB() {
	// Buscando credenciais das variáveis de ambiente (com fallbacks locais para o seu Ubuntu)
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "fiap_user")
	password := getEnv("DB_PASSWORD", "fiap_password")
	dbname := getEnv("DB_NAME", "fiap_x_db")

	// String de conexão padrão do Postgres (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Erro ao abrir configuração do banco: %v", err)
	}

	// Testa a conexão fisicamente para garantir que o banco está online e acessível
	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Banco de dados inacessível: %v", err)
	}

	fmt.Println("🗄️ Conexão com o PostgreSQL estabelecida com sucesso!")
}

// Função auxiliar para ler variáveis de ambiente ou usar um valor padrão
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// InsertVideo cria um registro usando seu UUID padrão e retorna o ID gerado como string
func InsertVideo(userID, originalName, storagePath string) (string, error) {
	var id string
	query := `
		INSERT INTO videos (user_id, original_name, storage_path, status) 
		VALUES ($1, $2, $3, 'PROCESSANDO') 
		RETURNING id`

	err := DB.QueryRow(query, userID, originalName, storagePath).Scan(&id)
	return id, err
}

// UpdateVideoSuccess atualiza o status para CONCLUIDO e salva os dados do ZIP usando o UUID
func UpdateVideoSuccess(id string, zipPath string, frameCount int) error {
	query := `
		UPDATE videos 
		SET status = 'CONCLUIDO', zip_path = $1, frame_count = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3`

	_, err := DB.Exec(query, zipPath, frameCount, id)
	return err
}

// UpdateVideoError marca o UUID como ERRO e registra a mensagem de falha do FFmpeg
func UpdateVideoError(id string, errorMessage string) error {
	query := `
		UPDATE videos 
		SET status = 'ERRO', error_message = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2`

	_, err := DB.Exec(query, errorMessage, id)
	return err
}

// GetVideosByWithUser busca o histórico de tarefas do usuário fictício
func GetVideosByUser(userID string) ([]models.Video, error) {
	query := `
		SELECT id, user_id, original_name, storage_path, COALESCE(zip_path, ''), frame_count, status, COALESCE(error_message, ''), created_at, updated_at 
		FROM videos 
		WHERE user_id = $1 
		ORDER BY created_at DESC`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		err := rows.Scan(&v.ID, &v.UserID, &v.OriginalName, &v.StoragePath, &v.ZipPath, &v.FrameCount, &v.Status, &v.ErrorMessage, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	return videos, nil
}
