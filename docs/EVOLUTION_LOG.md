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

---

