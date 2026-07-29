# Arquitetura Detalhada — FIAP X

Diagrama completo para referência interna. Cobre todos os componentes, fluxos de dados, protocolos, volumes compartilhados, segurança e observabilidade.

```mermaid
flowchart TD
    %% ─────────────────────────────────────────
    %% CLIENTE
    %% ─────────────────────────────────────────
    subgraph CLIENT["🖥️ Cliente (Browser)"]
        UI["HTML5 + JS Vanilla\n(Upload / Dashboard / Download)"]
    end

    %% ─────────────────────────────────────────
    %% GATEWAY
    %% ─────────────────────────────────────────
    subgraph GATEWAY["🔀 API Gateway (Nginx)"]
        direction TB
        RL["Rate Limiting\n10 req/s por IP\nburst 20 (geral) / 5 (upload)"]
        SH["Security Headers\nX-Frame-Options\nX-Content-Type-Options\nX-XSS-Protection\nReferrer-Policy"]
        PROXY["Proxy Reverso\nport 8080 → api-service:8080"]
        RL --> SH --> PROXY
    end

    %% ─────────────────────────────────────────
    %% API GO
    %% ─────────────────────────────────────────
    subgraph API["⚙️ API Backend (Go + Gin)"]
        direction TB
        MW["Auth Middleware\nJWT Validation"]

        subgraph ROUTES["Rotas HTTP"]
            R1["POST /login\n(público)"]
            R2["POST /upload\n🔒 JWT required"]
            R3["GET /api/videos\n🔒 JWT required"]
            R4["GET /download/:file\n🔒 JWT required"]
            R5["GET /metrics\n(Prometheus)"]
        end

        subgraph APP["Application Layer"]
            AS["AuthService\nbcrypt + JWT"]
            VS["VideoService\nOrchestration"]
        end

        subgraph ADAPTERS["Adapters"]
            RMQP["RabbitMQ Adapter\nPublishVideoPending()"]
            REPO["Postgres Repository\nInsertVideo()\nUpdateVideoError()\nGetVideosByUser()"]
        end

        MW --> ROUTES
        R1 --> AS
        R2 --> VS
        R3 --> REPO
        R4 --> REPO
        VS --> REPO
        VS --> RMQP
        AS --> REPO
    end

    %% ─────────────────────────────────────────
    %% MENSAGERIA
    %% ─────────────────────────────────────────
    subgraph MQ["🐇 RabbitMQ"]
        QUEUE["video_processing_queue\ndurable: true\npersistent messages"]
    end

    %% ─────────────────────────────────────────
    %% WORKER
    %% ─────────────────────────────────────────
    subgraph WORKER["🐍 Worker (Python)"]
        direction TB
        CONSUME["Consumer\nbasic_consume()\nprefetch=1"]

        subgraph WFLOW["Processing Flow"]
            W1["1. update → PROCESSANDO"]
            W2["2. process_video()\nFFmpeg: fps=1 → frames PNG"]
            W3["3. zipfile.ZipFile()\nframes_videoID.zip"]
            W4["4. update → CONCLUIDO\nzip_path + frame_count"]
            W5["5. ACK → RabbitMQ"]
            W1 --> W2 --> W3 --> W4 --> W5
        end

        subgraph WERR["Error Flow"]
            E1["update → ERRO\nerror_message"]
            E2["send_error_email()\nSMTP → Mailtrap"]
            E3["NACK → RabbitMQ\nrequeue: false"]
            E1 --> E2 --> E3
        end

        WMETRICS["Prometheus Metrics\nprocessed_total Counter\n:8000/metrics"]

        CONSUME --> WFLOW
        CONSUME --> WERR
        WFLOW --> WMETRICS
    end

    %% ─────────────────────────────────────────
    %% BANCO DE DADOS
    %% ─────────────────────────────────────────
    subgraph DB["🗄️ PostgreSQL 16"]
        direction TB
        TU["Table: users\nid UUID PK\nusername UNIQUE\npassword_hash bcrypt\nemail VARCHAR(255) nullable"]
        TV["Table: videos\nid UUID PK\nuser_id FK → users\noriginal_name\nstorage_path\nzip_path\nframe_count\nstatus ENUM\nerror_message"]
        ENUM["ENUM: video_status\nPENDENTE\nPROCESSANDO\nCONCLUIDO\nERRO"]
        TV --> ENUM
    end

    %% ─────────────────────────────────────────
    %% VOLUMES COMPARTILHADOS
    %% ─────────────────────────────────────────
    subgraph VOL["💾 Volumes Compartilhados (Docker/K8s PVC)"]
        VU["/app/uploads\nvídeos originais"]
        VO["/app/outputs\nZIPs gerados"]
        VT["/app/temp\nframes extraídos (PNG)"]
    end

    %% ─────────────────────────────────────────
    %% OBSERVABILIDADE
    %% ─────────────────────────────────────────
    subgraph OBS["📊 Observabilidade"]
        PROM["Prometheus\nscrape_interval: 5s\n:9090"]
        GRAF["Grafana\nDashboard\n:3000"]
        PROM --> GRAF
    end

    %% ─────────────────────────────────────────
    %% NOTIFICAÇÃO
    %% ─────────────────────────────────────────
    subgraph MAIL["📧 Mailtrap Sandbox"]
        INBOX["Inbox\nSMTP sandbox.smtp.mailtrap.io:2525\nalerta de falha de processamento"]
    end

    %% ─────────────────────────────────────────
    %% FLUXO PRINCIPAL
    %% ─────────────────────────────────────────
    UI -->|"HTTPS / HTTP"| GATEWAY
    GATEWAY -->|"proxy_pass"| API

    API -->|"INSERT videos (PENDENTE)"| DB
    API -->|"PublishMessage {video_id, video_path}"| MQ
    API -->|"salva arquivo"| VU

    MQ -->|"consume message"| WORKER
    WORKER -->|"UPDATE status"| DB
    WORKER -->|"lê vídeo"| VU
    WORKER -->|"escreve ZIP"| VO
    WORKER -->|"frames temporários"| VT

    API -->|"GET /download lê ZIP"| VO

    PROM -->|"scrape :8080/metrics"| API
    PROM -->|"scrape :8000/metrics"| WORKER

    WORKER -->|"SMTP"| MAIL

    %% ─────────────────────────────────────────
    %% ESTILOS
    %% ─────────────────────────────────────────
    classDef gateway fill:#f0ad4e,stroke:#d68910,color:#000
    classDef api     fill:#5dade2,stroke:#2e86c1,color:#000
    classDef worker  fill:#58d68d,stroke:#1e8449,color:#000
    classDef db      fill:#a569bd,stroke:#7d3c98,color:#fff
    classDef mq      fill:#f1948a,stroke:#c0392b,color:#000
    classDef obs     fill:#85c1e9,stroke:#2980b9,color:#000
    classDef vol     fill:#f9e79f,stroke:#d4ac0d,color:#000
    classDef mail    fill:#fdfefe,stroke:#aab7b8,color:#000
    classDef client  fill:#d5dbdb,stroke:#717d7e,color:#000

    class GATEWAY,RL,SH,PROXY gateway
    class API,MW,ROUTES,APP,ADAPTERS,R1,R2,R3,R4,R5,AS,VS,RMQP,REPO api
    class WORKER,CONSUME,WFLOW,WERR,WMETRICS,W1,W2,W3,W4,W5,E1,E2,E3 worker
    class DB,TU,TV,ENUM db
    class MQ,QUEUE mq
    class OBS,PROM,GRAF obs
    class VOL,VU,VO,VT vol
    class MAIL,INBOX mail
    class CLIENT,UI client
```

---

## Fluxo Completo de um Upload

```mermaid
sequenceDiagram
    actor User as Usuário
    participant GW as Gateway (Nginx)
    participant API as API Go
    participant PG as PostgreSQL
    participant MQ as RabbitMQ
    participant WK as Worker Python
    participant FS as Volume Compartilhado
    participant MT as Mailtrap

    User->>GW: POST /upload (JWT + arquivo)
    GW->>GW: Rate limit check (burst 5)
    GW->>GW: Injeta security headers
    GW->>API: proxy_pass

    API->>API: Valida JWT
    API->>API: Valida extensão do arquivo (.mp4, .avi...)
    API->>FS: Salva arquivo em /uploads
    API->>PG: INSERT videos (status=PENDENTE)
    API->>MQ: Publish {video_id, video_path}
    API-->>User: HTTP 202 Accepted

    MQ->>WK: Deliver message
    WK->>PG: UPDATE status=PROCESSANDO

    alt Processamento bem-sucedido
        WK->>FS: Lê vídeo de /uploads
        WK->>WK: FFmpeg extrai frames (1fps → PNGs)
        WK->>FS: Escreve frames em /temp
        WK->>FS: Cria ZIP em /outputs
        WK->>PG: UPDATE status=CONCLUIDO + zip_path + frame_count
        WK->>MQ: basic_ack
    else Erro no processamento
        WK->>PG: UPDATE status=ERRO + error_message
        WK->>MT: send_error_email (SMTP)
        WK->>MQ: basic_nack (requeue=false)
    end

    User->>GW: GET /api/videos (polling)
    GW->>API: proxy_pass
    API->>PG: SELECT videos WHERE user_id=...
    API-->>User: JSON com status atualizado

    User->>GW: GET /download/:filename
    GW->>API: proxy_pass
    API->>FS: Lê ZIP de /outputs
    API-->>User: ZIP file (Content-Disposition: attachment)
```
