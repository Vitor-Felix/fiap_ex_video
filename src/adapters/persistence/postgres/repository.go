package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"video-processor/domain/entities"

	_ "github.com/lib/pq"
)

// Repository é o adapter responsável pela persistência em PostgreSQL.
type Repository struct {
	db *sql.DB
}

// NewRepository cria a conexão com o PostgreSQL e devolve um repositório pronto para uso.
func NewRepository() (*Repository, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "fiap_user")
	password := getEnv("DB_PASSWORD", "fiap_password")
	dbname := getEnv("DB_NAME", "fiap_x_db")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	fmt.Println("🗄️ Conexão com PostgreSQL estabelecida com sucesso!")

	return &Repository{
		db: db,
	}, nil
}

// getEnv retorna uma variável de ambiente ou um valor padrão.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// InsertVideo cria um registro usando seu UUID padrão.
func (r *Repository) InsertVideo(userID, originalName, storagePath string) (string, error) {
	var id string

	query := `
		INSERT INTO videos (user_id, original_name, storage_path, status)
		VALUES ($1, $2, $3, 'PROCESSANDO')
		RETURNING id`

	err := r.db.QueryRow(query, userID, originalName, storagePath).Scan(&id)

	return id, err
}

// UpdateVideoSuccess atualiza o status para CONCLUIDO.
func (r *Repository) UpdateVideoSuccess(id, zipPath string, frameCount int) error {
	query := `
		UPDATE videos
		SET status = 'CONCLUIDO',
		    zip_path = $1,
		    frame_count = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`

	_, err := r.db.Exec(query, zipPath, frameCount, id)

	return err
}

// UpdateVideoError atualiza o status para ERRO.
func (r *Repository) UpdateVideoError(id, errorMessage string) error {
	query := `
		UPDATE videos
		SET status = 'ERRO',
		    error_message = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`

	_, err := r.db.Exec(query, errorMessage, id)

	return err
}

// GetVideosByUser retorna o histórico de processamento de um usuário.
func (r *Repository) GetVideosByUser(userID string) ([]entities.Video, error) {
	query := `
		SELECT id,
		       user_id,
		       original_name,
		       storage_path,
		       COALESCE(zip_path, ''),
		       frame_count,
		       status,
		       COALESCE(error_message, ''),
		       created_at,
		       updated_at
		FROM videos
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []entities.Video

	for rows.Next() {
		var v entities.Video

		if err := rows.Scan(
			&v.ID,
			&v.UserID,
			&v.OriginalName,
			&v.StoragePath,
			&v.ZipPath,
			&v.FrameCount,
			&v.Status,
			&v.ErrorMessage,
			&v.CreatedAt,
			&v.UpdatedAt,
		); err != nil {
			return nil, err
		}

		videos = append(videos, v)
	}

	return videos, nil
}
