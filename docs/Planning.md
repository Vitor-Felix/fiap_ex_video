# Planning

🏁 Milestone 1: Infraestrutura Base, Gitflow e CI/CD
Objetivo: Deixar o ambiente pronto, com banco de dados, fluxo de ramificação protegido e testes automatizados rodando antes de criar novas features.

[x] Issue 1.1: Executar e testar o projeto original localmente

Status: Concluído. Entendido o uso do ffmpeg e geração do ZIP.

[x] Issue 1.2: Modelagem e Criação do Banco de Dados (PostgreSQL)

Status: Concluído. Instância do Postgres configurada via Docker.

[x] Issue 1.3: Estrutura Hexagonal, Gitflow e Testes Iniciais

O que fazer: Reorganizar as pastas do Go para simular uma Arquitetura Hexagonal simples (domain/entities, ports/, adapters/). Criar as branches padrão master e develop. Escrever um teste unitário básico em Go (ex: validação do formato do arquivo ou tamanho permitido).

Foco de Avaliação FIAP: Arquitetura de Software (Clean Architecture/Hexagonal) e Qualidade de Software.

[ ] Issue 1.4: Pipeline de CI com GitHub Actions

O que fazer: Criar o arquivo .github/workflows/ci.yml. Configurar para disparar a cada Pull Request aberto contra a branch develop, executando automaticamente o comando go test ./.... Configurar proteção de branch no GitHub.

Foco de Avaliação FIAP: CI/CD e Governança de Código.

🔑 Milestone 2: Persistência e Segurança (Refatorando a API Go)
Objetivo: Garantir que o monolito inicial persista dados e autentique usuários de forma segura.

[x] Issue 2.1: Conectar o código Go ao PostgreSQL

Status: Concluído. O upload agora registra no banco e a listagem HTML lê os status reais (PENDENTE, PROCESSANDO, CONCLUIDO, ERRO) via polling assíncrono.

[ ] Issue 2.2: Implementar Autenticação de Usuário (JWT)

O que fazer: Criar uma tabela simples de usuários no Postgres. Desenvolver a rota de /login que valida as credenciais e devolve um token JWT. Proteger as rotas de upload e histórico para que o usuário autenticado acesse apenas os seus próprios registros.

Foco de Avaliação FIAP: Requisitos Funcionais de Segurança (Sistema protegido por usuário e senha).

🔀 Milestone 3: Gateway, Desacoplamento e Mensageria
Objetivo: Criar a porta de entrada única do ecossistema e separar o recebimento do vídeo do processamento real para suportar picos de carga.

[ ] Issue 3.1: Configurar o API Gateway

O que fazer: Adicionar um serviço de API Gateway leve (como Nginx ou KrakenD) no docker-compose.yml. Configurar o Gateway na porta padrão (ex: :8080) para interceptar e rotear o tráfego do frontend para a API interna em Go. O frontend deixa de falar diretamente com o microsserviço de backend.

Foco de Avaliação FIAP: Padrões de Arquitetura de Microsserviços e Ponto Único de Entrada.

[ ] Issue 3.2: Subir o RabbitMQ e Criar a Fila

O que fazer: Adicionar o container oficial do RabbitMQ ao ambiente Docker e declarar a fila de mensagens de processamento.

Foco de Avaliação FIAP: Mensageria e Resiliência do Sistema.

[ ] Issue 3.3: Transformar a API Go em Produtor (Producer)

O que fazer: Alterar o handler de upload do Go. Ao receber o arquivo, a API salva o estado inicial no Postgres como PENDENTE, publica o ID do vídeo na fila do RabbitMQ e retorna imediatamente o status HTTP 202 para o usuário. A API não processa mais o vídeo.

Foco de Avaliação FIAP: Processamento Assíncrono e Garantia de que nenhuma requisição será perdida em momentos de pico.

🐍 Milestone 4: O Novo Worker de Processamento (Segundo Microsserviço)
Objetivo: Criar o microsserviço especialista em background para consumir e processar as filas sob demanda.

[ ] Issue 4.1: Criar o Worker Consumidor (Python/Go)

O que fazer: Criar uma nova aplicação isolada (pode ser em Python devido à forte compatibilidade com manipulação de mídia). Esse script deve escutar a fila do RabbitMQ. Ao capturar uma mensagem, ele altera imediatamente o status do vídeo no Postgres para PROCESSANDO.

Foco de Avaliação FIAP: Desenvolvimento de Microsserviços e Desacoplamento.

[ ] Issue 4.2: Extração de Frames, Geração do ZIP e Alertas

O que fazer: Migrar a lógica do ffmpeg para o Worker. Ele realiza o processamento do vídeo, joga o .zip resultante em um diretório compartilhado (Volume Docker) e altera o status para CONCLUIDO. Adicionar uma integração simulada (ex: Mailtrap) para disparar um e-mail de alerta caso o bloco de captura caia em ERRO.

Foco de Avaliação FIAP: Processamento em Background paralelo e Notificação em caso de falhas.

🚀 Milestone 5: Orquestração (Kubernetes) e Observabilidade
Objetivo: Sair do nível local comum e preparar a aplicação para os padrões de produção em escala exigidos pelos avaliadores.

[ ] Issue 5.1: Mapear Manifestos Kubernetes (K8s)

O que fazer: Criar uma pasta k8s/ na raiz do projeto contendo os arquivos YAML básicos de Deployment, Service e ConfigMap para a API Go, o Worker e o banco. O foco aqui é provar documentalmente que o desenho da aplicação suporta escalabilidade horizontal via K8s.

Foco de Avaliação FIAP: Arquitetura e Infraestrutura autoescalável recomendada.

[ ] Issue 5.2: Monitoramento com Prometheus e Grafana

O que fazer: Subir os containers do Prometheus e do Grafana no compose local. Expor métricas nativas ou simples da aplicação (como contagem de uploads efetuados e erros de processamento) e plotar em um painel básico do Grafana para visualização na apresentação de vídeo.

Foco de Avaliação FIAP: Requisito Técnico de Monitoramento e Observabilidade.
