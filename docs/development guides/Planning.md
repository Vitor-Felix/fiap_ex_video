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

[x] Issue 1.4: Pipeline de CI com GitHub Actions

O que fazer: Criar o arquivo .github/workflows/ci.yml. Configurar para disparar a cada Pull Request aberto contra a branch develop, executando automaticamente o comando go test ./.... Configurar proteção de branch no GitHub.

Foco de Avaliação FIAP: CI/CD e Governança de Código.

🔑 Milestone 2: Persistência e Segurança (Refatorando a API Go)
Objetivo: Garantir que o monolito inicial persista dados e autentique usuários de forma segura.

[x] Issue 2.1: Conectar o código Go ao PostgreSQL

Status: Concluído. O upload agora registra no banco e a listagem HTML lê os status reais (PENDENTE, PROCESSANDO, CONCLUIDO, ERRO) via polling assíncrono.

[x] Issue 2.2: Implementar Autenticação de Usuário (JWT)

O que fazer: Criar uma tabela simples de usuários no Postgres. Desenvolver a rota de /login que valida as credenciais e devolve um token JWT. Proteger as rotas de upload e histórico para que o usuário autenticado acesse apenas os seus próprios registros.

Foco de Avaliação FIAP: Requisitos Funcionais de Segurança (Sistema protegido por usuário e senha).

🔀 Milestone 3: Gateway, Desacoplamento e Mensageria
Objetivo: Criar a porta de entrada única do ecossistema e separar o recebimento do vídeo do processamento real para suportar picos de carga.

[x] Issue 3.1: Configurar o API Gateway

O que fazer: Adicionar um serviço de API Gateway leve (como Nginx ou KrakenD) no docker-compose.yml. Configurar o Gateway na porta padrão (ex: :8080) para interceptar e rotear o tráfego do frontend para a API interna em Go. O frontend deixa de falar diretamente com o microsserviço de backend.

Foco de Avaliação FIAP: Padrões de Arquitetura de Microsserviços e Ponto Único de Entrada.

[x] Issue 3.2: Subir o RabbitMQ e Criar a Fila

O que fazer: Adicionar o container oficial do RabbitMQ ao ambiente Docker e declarar a fila de mensagens de processamento.

Foco de Avaliação FIAP: Mensageria e Resiliência do Sistema.

[x] Issue 3.3: Transformar a API Go em Produtor (Producer)

O que fazer: Alterar o handler de upload do Go. Ao receber o arquivo, a API salva o estado inicial no Postgres como PENDENTE, publica o ID do vídeo na fila do RabbitMQ e retorna imediatamente o status HTTP 202 para o usuário. A API não processa mais o vídeo.

Foco de Avaliação FIAP: Processamento Assíncrono e Garantia de que nenhuma requisição será perdida em momentos de pico.

🐍 Milestone 4: O Novo Worker de Processamento (Segundo Microsserviço)
Objetivo: Criar o microsserviço especialista em background para consumir e processar as filas sob demanda.

[x] Issue 4.1: Criar o Worker Consumidor (Python/Go)

O que fazer: Criar uma nova aplicação isolada (pode ser em Python devido à forte compatibilidade com manipulação de mídia). Esse script deve escutar a fila do RabbitMQ. Ao capturar uma mensagem, ele altera imediatamente o status do vídeo no Postgres para PROCESSANDO.

Foco de Avaliação FIAP: Desenvolvimento de Microsserviços e Desacoplamento.

[x] Issue 4.2: Extração de Frames, Geração do ZIP e Alertas

O que fazer: Migrar a lógica do ffmpeg para o Worker. Ele realiza o processamento do vídeo, joga o .zip resultante em um diretório compartilhado (Volume Docker) e altera o status para CONCLUIDO. Adicionar uma integração simulada (ex: Mailtrap) para disparar um e-mail de alerta caso o bloco de captura caia em ERRO.

Foco de Avaliação FIAP: Processamento em Background paralelo e Notificação em caso de falhas.

🚀 Milestone 5: Orquestração (Kubernetes) e Observabilidade
Objetivo: Sair do nível local comum e preparar a aplicação para os padrões de produção em escala exigidos pelos avaliadores.

[x] Issue 5.1: Mapear Manifestos Kubernetes (K8s)

O que fazer: Criar uma pasta k8s/ na raiz do projeto contendo os arquivos YAML básicos de Deployment, Service e ConfigMap para a API Go, o Worker e o banco. O foco aqui é provar documentalmente que o desenho da aplicação suporta escalabilidade horizontal via K8s.

Foco de Avaliação FIAP: Arquitetura e Infraestrutura autoescalável recomendada.

[x] Issue 5.2: Monitoramento com Prometheus e Grafana

O que fazer: Subir os containers do Prometheus e do Grafana no compose local. Expor métricas nativas ou simples da aplicação (como contagem de uploads efetuados e erros de processamento) e plotar em um painel básico do Grafana para visualização na apresentação de vídeo.

Foco de Avaliação FIAP: Requisito Técnico de Monitoramento e Observabilidade.
---
Issue 5.3 — Notificação de Falhas por E-mail (Mailtrap)
Objetivo

Atender ao requisito de negócio:

"Em caso de erro, um usuário pode ser notificado (e-mail ou outro meio de comunicação)."

Implementar uma notificação automática quando o Worker falhar no processamento de um vídeo.

Passo 1 — Criar uma conta no Mailtrap

Objetivo

Obter um servidor SMTP de testes sem necessidade de enviar e-mails reais.

O que fazer

Criar uma conta gratuita no Mailtrap.
Criar uma Inbox.
Copiar as credenciais SMTP:
Host
Porta
Usuário
Senha
Passo 2 — Configurar variáveis de ambiente

Objetivo

Evitar credenciais fixas no código.

O que fazer

Adicionar no Docker Compose/Kubernetes:

SMTP_HOST
SMTP_PORT
SMTP_USER
SMTP_PASSWORD
SMTP_FROM
Passo 3 — Criar um módulo de envio de e-mail

Objetivo

Isolar toda a lógica de SMTP.

O que fazer

Criar um novo arquivo:

worker/email_service.py

Responsável por:

abrir conexão SMTP
autenticar
montar a mensagem
enviar o e-mail
Passo 4 — Disparar o e-mail somente em caso de erro

Objetivo

Notificar apenas quando o processamento falhar.

O que fazer

No bloco de tratamento de exceção do Worker:

atualizar status para ERRO (já existente)
chamar
send_error_email(...)
Passo 5 — Conteúdo do e-mail

Objetivo

Facilitar a identificação da falha.

Conteúdo sugerido

vídeo
usuário
horário
mensagem da exceção

Exemplo:

Assunto:
Falha no processamento do vídeo

Corpo:

O processamento do vídeo falhou.

Vídeo:
video.mp4

Erro:
FFmpeg exited with code 1

Data:
2026-07-28 15:33
Passo 6 — Atualizar documentação

Objetivo

Facilitar a execução pelo avaliador.

O que fazer

Adicionar no README:

criação da conta Mailtrap
onde colocar as credenciais
como visualizar os e-mails enviados
Resultado esperado

Quando ocorrer qualquer exceção durante o processamento:

Upload
      ↓
RabbitMQ
      ↓
Worker
      ↓
Erro
      ↓
Atualiza banco
      ↓
Envia e-mail
---
Obs.: Entrei no mailtrap e já tenho um sandbox antigo que pode ser reaproveitado ou exlcuido (só permite 1 sandbox por vez no free plan). Mas já tenho a conta:

GoBarber
Total messages sent: 0
Integration
Auto Forward
Manual Forward
Access Rights
Sending
Send live emails with Email API/SMTP

Your plan includes: 4,000 emails, high deliverability, global support.

Try API/SMTP

SMTP


Email


API


POP3

Credentials
Reset Credentials
Host

sandbox.smtp.mailtrap.io
Port

25, 465, 587 or 2525
Username

c726054e2a7ca8
Password

****9a38
Auth

PLAIN, LOGIN and CRAM-MD5
TLS

Optional (STARTTLS on all ports)
Code Samples
cURL

PHP
Node.js

Python
C#

Ruby

Other
Copy
curl \
--ssl-reqd \
--url 'smtp://sandbox.smtp.mailtrap.io:2525' \
--user 'c726054e2a7ca8:****9a38' \
--mail-from from@example.com \
--mail-rcpt to@example.com \
--upload-file - <<EOF
From: Magic Elves <from@example.com>
To: Mailtrap Sandbox <to@example.com>
Subject: You are awesome!
Content-Type: multipart/alternative; boundary="boundary-string"

--boundary-string
Content-Type: text/plain; charset="utf-8"
Content-Transfer-Encoding: quoted-printable
Content-Disposition: inline

Congrats for sending test email with Mailtrap!

If you are viewing this email in your sandbox =E2=80=93 the integration works.
Now send your email using our SMTP server and integration of your choice!

Good luck! Hope it works.

--boundary-string
Content-Type: text/html; charset="utf-8"
Content-Transfer-Encoding: quoted-printable
Content-Disposition: inline

<!doctype html>
<html>
  <head>
    <meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3DUTF-8">
  </head>
  <body style=3D"font-family: sans-serif;">
    <div style=3D"display: block; margin: auto; max-width: 600px;" class=3D"main">
      <h1 style=3D"font-size: 18px; font-weight: bold; margin-top: 20px">Congrats for sending test email with Mailtrap!</h1>
      <p>If you are viewing this email in your sandbox =E2=80=93 the integration works.</p>
      <img alt=3D"Inspect with Tabs" src=3D"https://assets-examples.mailtrap.io/integration-examples/welcome.png" style=3D"width: 100%;">
      <p>Now send your email using our SMTP server and integration of your choice!</p>
      <p>Good luck! Hope it works.</p>
    </div>
    <!-- Example of invalid for email html/css, will be detected by Mailtrap: -->
    <style>
      .main { background-color: white; }
      a:hover { border-left-width: 1em; min-height: 2em; }
    </style>
  </body>
</html>

--boundary-string--
EOF
