# Arquitetura — FIAP X

```mermaid
flowchart LR
    %% ── Atores ──────────────────────────────────
    USER(["👤 Usuário"])

    %% ── Gateway ──────────────────────────────────
    subgraph GW["API Gateway"]
        NGINX["Nginx\n─────────────\nRate Limiting\nSecurity Headers\nProxy Reverso"]
    end

    %% ── Backend ──────────────────────────────────
    subgraph BACK["Backend"]
        API["API Go / Gin\n─────────────\nHTTP · JWT Auth\nUpload · Download\nPrometheus /metrics"]
    end

    %% ── Dados & Mensageria ───────────────────────
    subgraph INFRA["Infraestrutura"]
        PG[("PostgreSQL\n─────────\nusers · videos\nstatus ENUM")]
        MQ["RabbitMQ\n─────────────\nvideo_processing\n_queue"]
    end

    %% ── Worker ───────────────────────────────────
    subgraph PROC["Processamento"]
        WK["Worker Python\n─────────────\nFFmpeg · ZIP\nPrometheus :8000"]
    end

    %% ── Volumes ──────────────────────────────────
    subgraph STORE["Storage Compartilhado"]
        VOL[("Volume\n──────────\n/uploads\n/outputs\n/temp")]
    end

    %% ── Observabilidade ──────────────────────────
    subgraph OBS["Observabilidade"]
        PROM["Prometheus"]
        GRAF["Grafana"]
        PROM --> GRAF
    end

    %% ── Notificação ──────────────────────────────
    MAIL["📧 Mailtrap\n(alerta de falha)"]

    %% ── Fluxo principal ──────────────────────────
    USER -->|"HTTP"| GW
    GW -->|"proxy"| BACK
    BACK -->|"INSERT / SELECT"| PG
    BACK -->|"publish message"| MQ
    BACK <-->|"upload / download"| STORE

    MQ -->|"consume"| PROC
    PROC -->|"UPDATE status"| PG
    PROC <-->|"frames · zip"| STORE
    PROC -->|"SMTP"| MAIL

    %% ── Métricas ─────────────────────────────────
    PROM -.->|"scrape /metrics"| BACK
    PROM -.->|"scrape :8000"| PROC

    %% ── Estilos ──────────────────────────────────
    classDef gw      fill:#f0ad4e,stroke:#d68910,color:#000,rx:6
    classDef api     fill:#5dade2,stroke:#2e86c1,color:#000,rx:6
    classDef db      fill:#a569bd,stroke:#7d3c98,color:#fff,rx:6
    classDef mq      fill:#f1948a,stroke:#c0392b,color:#000,rx:6
    classDef worker  fill:#58d68d,stroke:#1e8449,color:#000,rx:6
    classDef obs     fill:#85c1e9,stroke:#2980b9,color:#000,rx:6
    classDef store   fill:#f9e79f,stroke:#d4ac0d,color:#000,rx:6
    classDef mail    fill:#fdfefe,stroke:#aab7b8,color:#000,rx:6
    classDef user    fill:#d5dbdb,stroke:#717d7e,color:#000,rx:6

    class NGINX gw
    class API api
    class PG db
    class MQ mq
    class WK worker
    class PROM,GRAF obs
    class VOL store
    class MAIL mail
    class USER user
```

---

| Componente | Tecnologia | Responsabilidade |
|:---|:---|:---|
| **API Gateway** | Nginx | Rate limiting, security headers, proxy reverso |
| **API Backend** | Go 1.25 + Gin | Autenticação JWT, upload, orquestração de fluxo |
| **Worker** | Python 3.11 + Pika | Consumo de fila, FFmpeg, geração de ZIP |
| **Mensageria** | RabbitMQ 4 | Desacoplamento assíncrono entre API e Worker |
| **Banco de Dados** | PostgreSQL 16 | Usuários, vídeos e ciclo de vida do processamento |
| **Storage** | Volume Docker / K8s PVC | Compartilhamento de arquivos entre API e Worker |
| **Observabilidade** | Prometheus + Grafana | Coleta e visualização de métricas em tempo real |
| **Notificação** | Mailtrap (SMTP) | Alerta por e-mail em caso de falha de processamento |
