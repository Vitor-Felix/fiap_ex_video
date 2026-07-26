# ==========================================
# ESTÁGIO 1: Compilação (Builder)
# ==========================================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copia e baixa dependências
COPY src/go.mod src/go.sum ./src/
WORKDIR /app/src
RUN go mod download

# Copia o código Go
COPY src/ /app/src/

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/video-processor main.go


# ==========================================
# ESTÁGIO 2: Imagem de Execução (Final)
# ==========================================
FROM alpine:3.19

RUN apk add --no-cache ffmpeg

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

ENV APP_ENV=production

RUN mkdir -p uploads outputs temp && \
    chown -R appuser:appgroup /app

# Copia o binário compilado do estágio 'builder'
COPY --from=builder /app/video-processor /app/video-processor

# Copia o frontend estático
COPY web/ /app/web/

USER appuser

EXPOSE 8080

CMD ["./video-processor"]
