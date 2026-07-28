# FIAP X - Processador de Vídeos

## 📌 Sobre o Projeto

Este repositório contém a versão evoluída do projeto **FIAP X**, um sistema em Go que recebe uploads de vídeos, realiza a extração de seus quadros (frames) utilizando o FFmpeg e disponibiliza os resultados compactados em formato `.zip` para download.  

O projeto agora conta com **persistência de dados relacional** para gerenciar o histórico das tarefas e os estados do processamento em tempo real.

---

## 🚀 Funcionalidades Atuais

- **Upload e Extração Automática:** Recebimento de arquivos de vídeo via API HTTP e extração de frames integrada usando FFmpeg.
- **Persistência de Estados:** Registro de logs de execução diretamente no banco de dados com estados dinâmicos (`PROCESSANDO`, `CONCLUIDO`, `ERRO`).
- **Dashboard de Histórico:** Interface Web nativa que lista todos os arquivos processados e realiza consultas assíncronas automáticas (polling) para atualizar o andamento das tarefas na tela.
- **Download Direto:** Disponibilização imediata do link de download do arquivo `.zip` gerado para os processos concluídos com sucesso.

---

## 🏗️ Arquitetura Atual

A aplicação está estruturada de forma modular em Go, centralizando o tráfego em um servidor HTTP e persistindo o ciclo de vida do vídeo em um container PostgreSQL dedicado.

### Componentes Ativos:

| Componente | Tecnologia | Papel Atual no Ecossistema |
| :--- | :--- | :--- |
| **Backend API** | Go (Golang) + Gin | Gerencia as rotas HTTP, validação de arquivos e orquestração do fluxo. |
| **Motor de Mídia** | FFmpeg | Utilitário de sistema executado em background pelo Go para extrair os frames. |
| **Interface (Frontend)**| HTML5 + JavaScript (Vanilla) | Interface de upload e painel de histórico dinâmico integrada via AJAX. |
| **Banco de Dados** | PostgreSQL | Armazena chaves UUID, metadados do vídeo e mensagens de erro para auditoria. |
| **Containerização** | Docker & Docker Compose | Isolamento dos ambientes da API, do ambiente runtime (Alpine + FFmpeg) e do Banco de Dados. |

---

## ✅ Requisitos Atendidos (Estado Atual)

### Funcionais
- [x] Extração de frames automática por vídeo enviado.
- [x] Listagem de status real dos vídeos (`PROCESSANDO`, `CONCLUIDO`, `ERRO`) refletida na interface do usuário.
- [x] Download do arquivo `.zip` com os frames extraídos através da listagem do painel.

### Técnicos
- [x] Persistência de dados ativa e estruturada (PostgreSQL com suporte a UUID e tipos ENUM).
- [x] Versionamento estruturado do código-fonte no GitHub.
- [x] Dockerização completa da aplicação para ambiente de desenvolvimento.

---

## 🧪 Stack Tecnológica Utilizada

- **Linguagem Principal:** Go (Golang 1.21+)
- **Framework Web:** Gin Gonic
- **Banco de Dados:** PostgreSQL 16
- **Processamento de Imagem:** FFmpeg (nativo no container)
- **Orquestração Local:** Docker & Docker Compose

---

## 🚀 Como Rodar a Aplicação

Todo o ecossistema atual (API Go, Utilitários de Mídia e PostgreSQL com o script de tabelas `init.sql`) é inicializado de forma integrada.

### 1. Subir o ambiente completo
Na raiz do projeto, execute o comando abaixo para construir a imagem customizada com FFmpeg e inicializar o banco de dados:
```bash
docker compose up --build
```

### 2. Acessar a aplicação
Abra o seu navegador e acesse o painel principal:

Plaintext
http://localhost:8080

### 3. Gerenciamento e Limpeza (Se necessário)
Caso realize alterações estruturais nos scripts SQL de inicialização do banco (db/init.sql), lembre-se de resetar os volumes persistidos do container rodando:

```
docker compose down -v
```

### 📚 Documentações
Os artefatos complementares e os logs de evolução do projeto podem ser consultados diretamente no diretório /docs.

---

## 🧪 Executando o CI Localmente

Os comandos abaixo replicam exatamente o que o GitHub Actions executa nos jobs `go-lint-and-test` e `python-lint-and-test`. Rode-os antes de abrir um PR para garantir que a pipeline passará.

### Pré-requisitos
- Go 1.21+
- golangci-lint (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)
- Python 3.12+ com `venv`

---

### Go — Lint

```bash
cd src
golangci-lint run --timeout=3m
```

### Go — Testes + Cobertura

```bash
cd src
go test ./... -v -race -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out
```

---

### Python (Worker) — Lint

```bash
# Cria e ativa um ambiente virtual com flake8
python3 -m venv .venv
source .venv/bin/activate
pip install flake8==7.1.1

# Roda o lint a partir do diretório worker
cd worker
flake8 . \
  --exclude=__pycache__,tests/__pycache__ \
  --max-line-length=120 \
  --extend-ignore=W503 \
  --statistics
```

### Python (Worker) — Testes

```bash
# Não requer banco de dados nem RabbitMQ — todas as dependências externas são mockadas
cd worker
python3 -m unittest discover -s tests -p "test_*.py" -v
```

---

## 📧 Notificação de Falhas por E-mail (Mailtrap)

O Worker Python envia automaticamente um e-mail de alerta sempre que o processamento de um vídeo falha. A integração usa o [Mailtrap Sandbox](https://mailtrap.io) como servidor SMTP de testes — nenhum e-mail real é enviado.

### Configuração

1. Crie uma conta gratuita em [mailtrap.io](https://mailtrap.io) e abra sua Inbox.

2. Copie as credenciais SMTP da seção **Credentials** da sua Inbox (Host, Port, Username, Password).

3. Crie o arquivo `.env` na raiz do projeto (já está no `.gitignore`, nunca será commitado):

```bash
# .env — credenciais locais, não commitar
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=<seu_usuario_mailtrap>
SMTP_PASSWORD=<sua_senha_mailtrap>
SMTP_FROM=worker@fiap-x.local
SMTP_TO=dev@fiap-x.local
```

4. Suba o ambiente normalmente — o Docker Compose lê o `.env` automaticamente:

```bash
docker compose up --build
```

### Como visualizar os e-mails

Após uma falha de processamento, acesse a Inbox no painel do Mailtrap em [mailtrap.io/inboxes](https://mailtrap.io/inboxes). O e-mail de alerta contém:

- ID e nome do arquivo de vídeo
- Mensagem de erro completa
- Data e hora da falha

### Kubernetes

Para o ambiente K8s, injete as credenciais reais no Secret via `kubectl patch` — sem alterar nenhum arquivo do repositório:

```bash
kubectl patch secret fiap-secrets \
  --type='json' \
  -p='[
    {"op":"replace","path":"/data/SMTP_USER","value":"'$(echo -n "<seu_usuario_mailtrap>" | base64)'"},
    {"op":"replace","path":"/data/SMTP_PASSWORD","value":"'$(echo -n "<sua_senha_mailtrap>" | base64)'"}
  ]'

kubectl rollout restart deployment/fiap-worker-deployment
```

As variáveis `SMTP_*` já estão mapeadas no Deployment do Worker (`k8s/05-worker.yaml`).

---

## ☸️ Guia de Execução Local (Kubernetes via Minikube)
Para a execução completa no Minikube, incluindo detalhes de pré-requisitos, build das imagens, criação de ConfigMaps, deployment dos manifestos, acesso à aplicação, PostgreSQL, Prometheus e seed de usuários, consulte o guia detalhado em [k8s/README.md](k8s/README.md).

Resumo rápido:

```bash
minikube start

docker build -t fiap-api:v2 -f Dockerfile .
docker build -t fiap-worker:v3 -f worker/Dockerfile ./worker

minikube image load fiap-api:v2
minikube image load fiap-worker:v3

kubectl create configmap nginx-config --from-file=nginx/nginx.conf
kubectl create configmap postgres-init-script --from-file=db/init.sql

kubectl apply -f k8s/
```

Após aplicar os manifestos, injete as credenciais SMTP reais no Secret (não altera nenhum arquivo do repositório):

```bash
kubectl patch secret fiap-secrets \
  --type='json' \
  -p='[
    {"op":"replace","path":"/data/SMTP_USER","value":"'$(echo -n "<seu_usuario_mailtrap>" | base64)'"},
    {"op":"replace","path":"/data/SMTP_PASSWORD","value":"'$(echo -n "<sua_senha_mailtrap>" | base64)'"}
  ]'

kubectl rollout restart deployment/fiap-worker-deployment
```

Aguarde todos os pods ficarem `Running`:

```bash
kubectl get pods
```

Acesse a aplicação:

```bash
kubectl port-forward svc/gateway-service 8080:8080
```

```text
http://localhost:8080
```

Para inspecionar o banco via DBeaver ou psql:

```bash
kubectl port-forward svc/postgres-service 5432:5432
```

```text
Host: localhost  |  Port: 5432  |  Database: fiap_x_db  |  User: fiap_user
```

Se precisar criar usuários iniciais para login, siga o exemplo descrito no guia detalhado em [k8s/README.md](k8s/README.md).
