-- db/init.sql

-- Criação de um tipo ENUM para garantir a integridade dos estados no banco
CREATE TYPE video_status AS ENUM ('PENDENTE', 'PROCESSANDO', 'CONCLUIDO', 'ERRO');

CREATE TABLE IF NOT EXISTS videos (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id VARCHAR(255) NOT NULL,
	original_name VARCHAR(255) NOT NULL,
	storage_path VARCHAR(512) NOT NULL,
	zip_path VARCHAR(512),
	status video_status DEFAULT 'PENDENTE',
	error_message TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índice para acelerar a busca de vídeos por usuário
CREATE INDEX IF NOT EXISTS idx_videos_user_id ON videos(user_id);
