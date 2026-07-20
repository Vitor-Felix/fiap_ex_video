# EVOLUTION_LOG

Este arquivo registra o histórico exato de modificações na base de código, refatorações estruturais e novos artefatos técnicos inseridos no projeto para guiar os próximos agentes de IA sem perda de contexto.

🕒 [FASE 1] - Setup da Infraestrutura & Banco de Dados (Issues 1.1 e 1.2)
📁 Modificações no Repositório:
docker-compose.yml (Novo): Criação do ambiente isolado contendo o container oficial do PostgreSQL 16, mapeamento de portas (5432:5432), variáveis de ambiente para credenciais administrativas e montagem de volume local persistente (pgdata).

db/init.sql (Novo): Script de migração inicial executado automaticamente na subida do banco.

Habilitação da extensão "uuid-ossp" para geração automática de identificadores únicos.

Criação do tipo ENUM video_status contendo os estados: PENDENTE, PROCESSANDO, CONCLUIDO, ERRO.

Criação da tabela videos contendo os campos: id (UUID PRIMARY KEY), filename (VARCHAR), status (video_status), error_message (TEXT), created_at (TIMESTAMP).

🕒 [FASE 2] - Quebra do Monólito & Integração de Persistência (Issue 2.1 Concluída)
🏗️ Mudança Arquitetural (Fatiamento do main.go):
O projeto original consistia em um arquivo único main.go inflado, acumulando rotas, handlers HTTP, comandos de sistema (FFmpeg) e lógica de view. O código foi desacoplado e distribuído na seguinte estrutura física atual:

main.go (Refatorado): Atua estritamente como ponto de entrada (entrypoint). Inicializa a conexão com o banco de dados via driver lib/pq, instancia o roteador do framework Gin, registra os handlers e levanta o servidor HTTP na porta :8080.

db.go (Novo): Centraliza o gerenciamento do pool de conexões do PostgreSQL através de uma variável global DB *sql.DB. Contém a função de inicialização que valida o handshake com o banco.

handlers.go (Novo): Centraliza as funções controladoras (handlers) do Gin:

UploadHandler: Intercepta o arquivo multipart, gera um UUID, persiste a linha inicial no banco com status PENDENTE, invoca o fluxo de processamento e lida com a resposta HTTP.

StatusHandler: Endpoint de leitura que executa um SELECT filtrando metadados e status reais para alimentar a tela.

DownloadHandler: Endpoint que serve o arquivo .zip final do disco para o cliente.

processor.go (Novo): Encapsula toda a lógica pesada de I/O e execução do motor de mídia. Executa em background o comando os/exec disparando o utilitário binário FFmpeg instalado na imagem Alpine Linux, gerencia as pastas temporárias de frames, empacota o resultado final em .zip e atualiza as transições de estado (PROCESSANDO -> CONCLUIDO / ERRO) no banco de dados.

💻 Integração de Polling no Frontend:
index.html / app.js (Modificados): A interface do usuário foi limpa de lógicas mockadas. Agora, ao finalizar um upload, o JavaScript nativo inicia um ciclo de polling assíncrono via setInterval consumindo o endpoint /api/status. A tabela de histórico e os botões de download são renderizados dinamicamente via AJAX com base nos retornos reais do PostgreSQL.

🚨 Contexto Atual para o Próximo Agente:
O projeto compila com sucesso, roda de forma integrada via Docker Compose e executa todo o fluxo de persistência e conversão.

O próximo passo obrigatório na esteira é iniciar a Issue 1.3: Estrutura Hexagonal, Gitflow e Testes Iniciais, movendo essa separação preliminar de arquivos (db.go, handlers.go, processor.go) para pastas formais de Domínio, Portas e Adaptadores, sem quebrar os contratos estabelecidos.

🕒 [FASE 3] - Refatoração para Arquitetura Hexagonal & Testes Nativos (Issue 1.3 Concluída)
🏗️ Mudança Arquitetural (Ports and Adapters):
O código foi reestruturado para o padrão de Arquitetura Hexagonal (Ports and Adapters), garantindo a Inversão de Dependência e o isolamento total da Regra de Negócio. Arquivos genéricos anteriores (db.go, handlers.go, processor.go) foram removidos e distribuídos na seguinte topologia:

- domain/entities/ (Core): Contém o modelo anêmico de domínio (video.go).
- ports/outbound/ (Contratos): Interfaces `VideoProcessor` e `VideoRepository` que ditam como a aplicação se comunica com o mundo externo.
- application/ (Use Cases): O coração da aplicação (`video_service.go`). Orquestra a regra de negócio dependendo apenas das portas (interfaces), sem conhecimento de infraestrutura.
- adapters/persistence/postgres/: Implementação concreta do repositório (`repository.go`).
- adapters/ffmpeg/: Implementação concreta do processamento de vídeo e zip (`processor.go`, `zip.go`).
- adapters/web/: Camada de entrada HTTP utilizando Gin (`handler.go`, `upload.go`, `download.go`, `status.go`, `video_list.go`, `views.go`). Não possui acesso direto ao banco, comunicando-se exclusivamente via `VideoService`.
- dto/ e utils/: Objetos de transferência de dados e funções utilitárias isoladas (`files.go`).

🧪 Implementação de Testes Unitários Nativos:
Foi implementada uma suíte de testes unitários sem dependência de frameworks externos (exclusivamente pacotes nativos `testing` e `net/http/httptest`), focada no Core e nas entradas.
- Testes de Aplicação (`video_service_test.go`): Utilização do padrão "Fakes" (structs locais mockando as interfaces do Repositório e do Processador) para testar os caminhos de sucesso e falha (erros de banco) em milissegundos, isolando o Service do PostgreSQL e do FFmpeg reais.
- Testes de Utilitários (`files_test.go`): Validação de criação de diretórios e manipulação de File System utilizando `os` e arquivos temporários.
- Testes da Camada Web (`download_test.go`, `status_test.go`, `video_list_test.go`): Utilização do `httptest.NewRecorder()` e `gin.TestMode` para simular requisições HTTP, validar status codes (200 OK, 404 Not Found), headers e payloads JSON sem subir um servidor real. Injeção de repositórios Fake diretamente nos Handlers.

🚨 Contexto Atual para o Próximo Agente:
A base de código atual possui separação de responsabilidades clara e cobertura de testes sólida nas camadas vitais (Application, Web e Utils). O código já foi commitado nas branches `master` e `develop`.

O próximo passo obrigatório é iniciar a Issue 1.4: Configurar o Pipeline de CI/CD com GitHub Actions, criando os workflows `.yaml` para executar o linting e a suíte de testes (`go test -cover ./...`) automaticamente a cada Push/Pull Request na branch `develop`.
