-- db/init.sql

-- Habilita funções de geração de UUID
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Criação do ENUM responsável pelos estados do processamento
CREATE TYPE video_status AS ENUM (
    'PENDENTE',
    'PROCESSANDO',
    'CONCLUIDO',
    'ERRO'
);

-- =====================================================
-- Tabela de usuários
-- =====================================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Tabela de vídeos
-- =====================================================
CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Dono do vídeo
    user_id UUID NOT NULL,

    original_name VARCHAR(255) NOT NULL,
    storage_path VARCHAR(512) NOT NULL,
    zip_path VARCHAR(512),
    frame_count INT DEFAULT 0,

    status video_status DEFAULT 'PENDENTE',
    error_message TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_videos_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Índice para acelerar buscas por usuário
CREATE INDEX IF NOT EXISTS idx_videos_user_id
ON videos(user_id);
