# EVOLUTION_LOG

Este arquivo registra o histórico exato de modificações na base de código, refatorações estruturais e novos artefatos técnicos inseridos no projeto para guiar os próximos agentes de IA sem perda de contexto.

🕒 [FASE 1] - Setup da Infraestrutura & Banco de Dados (Issues 1.1 e 1.2)
📁 Modificações no Repositório:
docker-compose.yml (Novo): Criação do ambiente isolado contendo o container oficial do PostgreSQL 16, mapeamento de portas (5432:5432), variáveis de ambiente para credenciais administrativas e montagem de volume local persistente (pgdata).

db/init.sql (Novo): Script de migração inicial executado automaticamente na subida do banco.
- Habilitação da extensão "uuid-ossp" para geração automática de identificadores únicos.
- Criação do tipo ENUM video_status contendo os estados: PENDENTE, PROCESSANDO, CONCLUIDO, ERRO.
- Criação da tabela videos contendo os campos: id (UUID PRIMARY KEY), filename (VARCHAR), status (video_status), error_message (TEXT), created_at (TIMESTAMP).

🕒 [FASE 2] - Quebra do Monólito & Integração de Persistência (Issue 2.1 Concluída)
🏗️ Mudança Arquitetural (Fatiamento do main.go):
O projeto original consistia em um arquivo único main.go inflado, acumulando rotas, handlers HTTP, comandos de sistema (FFmpeg) e lógica de view. O código foi desacoplado e distribuído em módulos dedicados.

💻 Integração de Polling no Frontend:
index.html / app.js (Modificados): A interface do usuário foi limpa de lógicas mockadas. O JavaScript nativo iniciou um ciclo de polling assíncrono via `setInterval` consumindo retornos reais do banco de dados.

🕒 [FASE 3] - Refatoração para Arquitetura Hexagonal & Testes Nativos (Issue 1.3 Concluída)
🏗️ Mudança Arquitetural (Ports and Adapters):
O código foi reestruturado para o padrão de Arquitetura Hexagonal (Ports and Adapters), garantindo a Inversão de Dependência e o isolamento total da Regra de Negócio:
- `domain/entities/` (Core): Contém o modelo anêmico de domínio (`video.go`, `user.go`).
- `ports/outbound/` (Contratos): Interfaces `VideoProcessor`, `VideoRepository` e `UserRepository`.
- `application/` (Use Cases): O coração da aplicação (`video_service.go`, `auth_service.go`). Orquestra a regra de negócio dependendo apenas das portas.
- `adapters/persistence/postgres/`: Implementação concreta do repositório (`repository.go`).
- `adapters/ffmpeg/`: Implementação concreta do processamento de vídeo e zip (`processor.go`, `zip.go`).
- `adapters/web/`: Camada de entrada HTTP utilizando Gin (`handler.go`, `upload.go`, `download.go`, `video_list.go`, `views.go`, `auth_middleware.go`).
- `dto/` e `utils/`: DTOs e utilitários isolados (`files.go`).

🧪 Implementação de Testes Unitários Nativos:
Suíte de testes sem frameworks externos (`testing` e `net/http/httptest`) focada em Application, Web e Utils usando Fakes/Mocks manuais.

🕒 [FASE 4] - Pipeline de CI/CD & Proteção de Branches (Issue 1.4 Concluída)
⚙️ Configuração do GitHub Actions:
- Pipeline `.github/workflows/ci.yml` configurado com Linting (`golangci-lint` v1.64.8) e `go test -race -cover`.
- Proteção de branches `main` e `develop` via Gitflow.

🕒 [FASE 5] - Autenticação JWT & Interface SPA Desacoplada (Issue 2.2 Concluída)

🔐 Modelagem de Dados & Autenticação (Backend)
- Tabela `users` criada com extensão `pgcrypto` para geração de `UUID`s.
- Tabela `videos` refatorada: `user_id` alterado de `VARCHAR` para `UUID` com `FOREIGN KEY` (`ON DELETE CASCADE`) apontando para `users(id)`.
- Senhas protegidas com `bcrypt` (`golang.org/x/crypto/bcrypt`).
- Autenticação stateless via JWT (`github.com/golang-jwt/jwt/v5`).
- `AuthMiddleware`: Intercepta o header `Authorization: Bearer <token>`, valida a assinatura e disponibiliza `user_id` e `username` no contexto do Gin (`c.Set(...)`).
- Todas as rotas de negócio (`/upload`, `/download/:filename`, `/api/videos`) protegidas por JWT. Cada usuário visualiza e gerencia apenas seus próprios arquivos.

🎨 Interface de Usuário (Frontend Single Page Application)
- Desacoplamento da renderização do HTML: Frontend migrado para SPA Vanilla JS e HTML5 sem frameworks pesados em `web/static/index.html` e `web/static/app.js`.
- Gerenciamento de Sessão: Token JWT mantido no `localStorage` do navegador.
- Experiência de Usuário:
  - Form de Login com transição dinâmica para a Dashboard do Usuário.
  - Botão de Logout com limpeza do token e reset da interface.
  - Polling dinâmico via `setInterval` chamando `/api/videos` com header `Authorization` injetado.
  - Botão de Download consumindo a API autenticada via `fetch` (Blob) para forçar o salvamento do arquivo ZIP.

📂 Gestão de Caminhos e Servimento de Arquivos Estáticos (Ajustes Críticos de Infra)
- Resolução de Erros de Caminho Relativo (404/401):
  - Inclusão do `utils.BasePath` para apontar corretamente para a raiz do repositório a partir de qualquer diretório de execução do terminal.
  - `views.go`: Ajustado para servir `web/static/index.html` e os assets CSS/JS em `/static` de forma consistente.
- Estrutura de Diretórios Final na Raiz do Repositório:
  ```text
  fiap_ex_video/
  ├── db/                # Scripts SQL
  ├── docs/              # Documentações
  ├── src/               # Código-fonte Go (Hexagonal)
  ├── uploads/           # Arquivos de upload temporários
  ├── outputs/           # Arquivos ZIP finais
  └── web/               # Assets do Frontend SPA
      └── static/
          ├── app.js
          └── index.html
  ```

🕒 [FASE 6] - API Gateway & Isolamento de Rede (Issue 3.1 Concluída)

🚪 API Gateway (Nginx)

Introdução do serviço gateway.
Porta pública concentrada em localhost:8080.
Backend Go deixou de expor portas diretamente ao host.
Comunicação interna realizada através da rede Docker.

🐳 Container da API

Atualização do Dockerfile para golang:1.25-alpine.
Inclusão dos assets do frontend na imagem (COPY web/ /app/web/).

🕒 [FASE 7] - Infraestrutura de Mensageria (Issue 3.2 Concluída)

📨 Introdução do RabbitMQ

A arquitetura foi preparada para evoluir de processamento síncrono para processamento assíncrono através de mensageria.

docker-compose.yml (Modificado)

Novo serviço:

rabbitmq:4-management

Configurações adicionadas:

Interface administrativa (15672)
Porta AMQP (5672)
Usuário e senha padrão para desenvolvimento
Volume persistente (rabbitmq_data)
Health Check utilizando rabbitmq-diagnostics
Hostname interno rabbitmq

A API também passou a aguardar o RabbitMQ ficar saudável (depends_on) antes da inicialização.

🏗️ Mudança Arquitetural

Até esta etapa, a arquitetura ficou organizada da seguinte forma:

                Usuário
                    │
                    ▼
             Nginx Gateway
                    │
                    ▼
                API Go
               /      \
              ▼        ▼
        PostgreSQL   RabbitMQ

O RabbitMQ foi incorporado apenas como componente de infraestrutura.

Nenhuma regra de negócio foi alterada nesta Issue.

A API continua executando o FFmpeg localmente e o fluxo permanece síncrono.

📚 Conceitos consolidados durante esta etapa

Foi estabelecido o desenho arquitetural que será utilizado nas próximas Issues:

A API Go atuará como Producer.
O Worker atuará como Consumer.
O RabbitMQ será responsável apenas pelo transporte de mensagens.
O PostgreSQL continuará sendo a fonte de verdade do estado da aplicação.
A mensagem enviada ao RabbitMQ conterá apenas os dados necessários para localizar o trabalho (ex.: video_id), nunca o vídeo em si.

Também foi discutido o papel do RabbitMQ em arquiteturas distribuídas:

desacoplamento entre serviços;
comunicação assíncrona;
possibilidade de múltiplos Workers;
escalabilidade horizontal;
mecanismo de confirmação (ACK) para evitar perda de mensagens.

💡 Instruções Importantes para a Próxima IA / Desenvolvedor

Ambiente atualizado:

docker compose up --build -d

Serviços esperados:

fiap_gateway
fiap_api
fiap_postgres
fiap_rabbitmq

Interface administrativa do RabbitMQ:

http://localhost:15672

Credenciais:

Usuário: fiap
Senha: fiap

Testes continuam sendo executados normalmente:

cd src
go test ./...
Próxima Etapa Oficial
Issue 3.3 — Transformar a API Go em Producer

Objetivo:

Remover o processamento síncrono da API.
Após o upload:
persistir o vídeo como PENDENTE;
publicar uma mensagem (video_id) no RabbitMQ;
retornar imediatamente HTTP 202 Accepted.
O processamento via FFmpeg deixará de ocorrer na API e será migrado para o Worker na Milestone 4.

Observação importante para a próxima IA: nesta etapa ainda não existe Consumer. O foco será exclusivamente integrar a API Go ao RabbitMQ, declarar a fila (QueueDeclare) e publicar mensagens (Publish) utilizando um cliente AMQP. O Worker será implementado apenas na Issue 4.1.

🕒 [FASE 8] - API Go como Producer & Desacoplamento Assíncrono (Issue 3.3 Concluída)

📨 Integração do RabbitMQ na API Go (Producer)
- Criação da interface outbound `MessageBroker` em `src/ports/outbound/message_broker.go` respeitando a Arquitetura Hexagonal.
- Criação do adaptador `RabbitMQAdapter` em `src/adapters/messaging/rabbitmq.go` utilizando o driver oficial `github.com/rabbitmq/amqp091-go`.
  - Configuração de fila durável (`video_processing_queue`) com mensagens persistentes.
  - Envio de payload JSON com `video_id` e `video_path` via `PublishWithContext` com timeout de segurança.

🔄 Alteração de Fluxo da Regra de Negócio (`video_service.go` & `upload.go`)
- Remoção da chamada síncrona do FFmpeg no fluxo HTTP.
- O handler `/upload` armazena o vídeo físico, grava o registro no banco como `PENDENTE`, enfileira a mensagem no RabbitMQ e responde imediatamente com **HTTP 202 Accepted**.
- O arquivo de vídeo físico em `/app/uploads` é mantido intacto no disco compartilhado para consumo posterior pelo Worker.

🐳 Ajuste de Volumes e Permissões no Docker Compose
- Transição de Bind Mounts locais para **Named Volumes** compartilhados (`shared_uploads`, `shared_outputs`, `shared_temp`) no `docker-compose.yml`.
- Resolução definitiva de conflitos de permissão (`permission denied`) mantendo compatibilidade nativa e execução limpa sem intervenção manual do usuário/avaliador.

🧪 Atualização da Suíte de Testes Nativos
- Refatoração dos mocks em `video_service_test.go`: substituição do `fakeProcessor` pelo `fakeBroker`.
- Teste unitário validando a publicação na fila em cenário de sucesso e o aborto da publicação em caso de erro de persistência.
- Cobertura e linting aprovados (`go test ./...` e `golangci-lint run ./...`).

🕒 [FASE 9] - Introdução do Worker Consumidor Poliglota em Python (Issue 4.1 Concluída)

🐍 Novo Microsserviço Worker (Python)

Criação do diretório /worker isolado na raiz do repositório contendo o novo microsserviço especialista consumidor de filas.

worker/requirements.txt: Dependências leves com pika==1.3.2 (AMQP) e psycopg2-binary==2.9.9 (Driver PostgreSQL).

worker/db.py: Módulo isolado de persistência para atualização de estado dos vídeos no banco relacional.

worker/main.py: Script principal de escuta contínua no RabbitMQ.

Implementação de retry resilience para estabilidade na inicialização do container.

Consumo com confirmação manual (basic_ack em caso de sucesso e basic_nack em falhas).

Configuração de Fair Dispatch (prefetch_count=1).

Transição imediata do status do vídeo de PENDENTE para PROCESSANDO via Query SQL direta no Postgres.

worker/Dockerfile: Imagem base python:3.11-slim otimizada com instalação de dependências C para suporte ao psycopg2.

🐳 Atualização da Orquestração no docker-compose.yml

Inclusão do serviço fiap_worker apontando para o contexto ./worker.

Mapeamento correto de variáveis de ambiente apontando para o container database (Postgres fiap_x_db / fiap_user) e rabbitmq (AMQP porta 5672).

Montagem dos Named Volumes compartilhados (shared_uploads, shared_outputs, shared_temp).

💡 Instruções Importantes para a Próxima IA / Desenvolvedor
Ambiente Atualizado e Operacional:

Bash
docker compose up --build -d
Serviços Esperados e Ativos:

fiap_gateway (Nginx na porta :8080)

fiap_api (Backend Go / Producer)

fiap_postgres (PostgreSQL 16 / Banco fiap_x_db)

fiap_rabbitmq (RabbitMQ 4 / Fila video_processing_queue)

fiap_worker (Worker Python / Consumer)

📌 Estado Atual do Fluxo Assíncrono:
O Usuário faz o upload na SPA.

A API Go valida, armazena o vídeo em /app/uploads, cria o registro como PENDENTE no Postgres e envia o JSON {"video_id": "...", "video_path": "..."} para o RabbitMQ.

A API Go responde imediatamente HTTP 202 Accepted.

O Worker Python captura a mensagem, abre conexão com o Postgres e altera o status do vídeo para PROCESSANDO.

O Worker envia o ACK para o RabbitMQ.

A interface SPA atualiza a badge via Polling para PROCESSANDO.

🚀 Próxima Etapa Oficial:
Milestone 4 — Issue 4.2: Extração de Frames via FFmpeg, Geração do ZIP e Finalização/Erro

Objetivo:

Adicionar a binary do ffmpeg no Dockerfile do Worker Python (ou via pacote de sistema apt-get install -y ffmpeg).

No script Python do Worker (main.py / novo módulo processor.py):

Executar o ffmpeg extraindo 1 frame por segundo (-vf fps=1) do arquivo em /app/uploads.

Salvar os frames temporariamente no volume /app/temp.

Compactar os frames extraídos em um arquivo .zip na pasta /app/outputs.

Em caso de SUCESSO: atualizar no Postgres para status = 'CONCLUIDO'.

Em caso de ERRO: capturar a exceção e atualizar para status = 'ERRO' gravando a mensagem em error_message.

Opcional/Diferencial FIAP: Adicionar um mock/módulo de notificação por e-mail em caso de falha.

Atenção para a Próxima IA:

Toda a lógica do FFmpeg que existia no Go em src/adapters/ffmpeg/ deve agora ser implementada em Python dentro da pasta /worker.

Lembre-se de instalar o ffmpeg no worker/Dockerfile.

Os caminhos de volume compartilhados já estão configurados no docker-compose.yml (/app/uploads, /app/outputs, /app/temp).


---

<FollowUp label="Quer iniciar agora a Issue 4.2 para instalar o FFmpeg e fazer a extração de frames e ZIP no Python?" query="Vamos começar a Issue 4.2! Como adicionamos o FFmpeg e a lógica de geração de ZIP no nosso Worker Python?"/>
