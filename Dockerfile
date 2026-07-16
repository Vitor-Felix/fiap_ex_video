# ==========================================
# ESTÁGIO 1: Compilação (Builder)
# ==========================================
FROM golang:1.21-alpine AS builder

# Instala ferramentas necessárias para compilação (se houver dependência C)
RUN apk add --no-cache git

WORKDIR /app

# Copia apenas os arquivos de dependência primeiro (otimiza o cache de camadas do Docker)
COPY src/go.mod src/go.sum ./src/
WORKDIR /app/src
RUN go mod download

# Copia o resto do código-fonte Go
COPY src/ /app/src/

# Compila o binário estático do Go otimizado para produção:
# - ldflags="-s -w": Remove tabelas de símbolos de debug para reduzir o tamanho do binário
# - CGO_ENABLED=0: Desabilita dependências dinâmicas de C, gerando um binário 100% portável
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/video-processor main.go


# ==========================================
# ESTÁGIO 2: Imagem de Execução (Final)
# ==========================================
FROM alpine:3.19

# Instala o ffmpeg necessário para rodar o comando em runtime
RUN apk add --no-cache ffmpeg

# Cria um usuário de sistema não-privilegiado por segurança
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

ENV APP_ENV=production

# Cria as pastas de persistência na raiz e define a permissão para o nosso usuário
RUN mkdir -p uploads outputs temp && \
    chown -R appuser:appgroup /app

# Copia apenas o executável compilado do estágio anterior (Builder)
COPY --from=builder /app/video-processor /app/video-processor

# Altera o contexto de execução para o usuário seguro
USER appuser

EXPOSE 8080

# Executa diretamente o binário compilado, sem go run!
CMD ["./video-processor"]
