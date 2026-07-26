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

💡 Instruções Importantes para a Próxima IA / Desenvolvedor:

Ponto de Execução do Servidor: A aplicação deve sempre ser iniciada a partir da raiz do repositório para evitar inconsistências nos caminhos de arquivos estáticos e diretórios do sistema:

Bash
cd ~/vitorfelix-git/fiap/fiap_ex_video
go run src/main.go
Execução dos Testes: Todos os testes continuam passando. Execute via:

Bash
cd src && go test ./...
Próxima Etapa do Projeto: O sistema está 100% funcional de ponta a ponta, seguro e autenticado. O próximo passo recomendado é Issue 2.3: Processamento Assíncrono com Worker Pool / Goroutines, focando em otimizar a fila de conversão de vídeos para que múltiplos uploads não travem a thread HTTP principal.
