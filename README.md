# FIAP EX Video - Processamento de Vídeos em Escala

## 📌 Sobre o Projeto

Este repositório contém a versão evoluída do projeto **FIAP X**, um sistema de processamento de vídeos que extrai quadros (frames) e os disponibiliza para download em formato `.zip`.  
A nova arquitetura foi projetada para atender requisitos de **alta disponibilidade**, **escalabilidade**, **segurança** e **resiliência**, conforme demanda dos investidores.

---

## 🚀 Funcionalidades

- Processamento simultâneo de múltiplos vídeos
- Fila de mensagens para não perder requisições em momentos de pico
- Autenticação e autorização (usuário e senha)
- Listagem de status dos vídeos enviados por usuário
- Notificação por e-mail em caso de erro no processamento
- Download do arquivo `.zip` com os frames extraídos

---

## 🏗️ Arquitetura Proposta

![Arquitetura do Sistema](./docs/arquitetura.png) *(adicione a imagem em /docs)*

### Componentes principais:

| Componente       | Tecnologia                     |
|------------------|--------------------------------|
| API Gateway      | Spring Cloud Gateway / Nginx   |
| Autenticação     | Keycloak / JWT                 |
| Microsserviço de upload | Spring Boot (ou Node.js) |
| Microsserviço de processamento | Python (OpenCV) + Flask/FastAPI |
| Mensageria       | RabbitMQ                       |
| Banco de dados   | PostgreSQL                     |
| Cache            | Redis                          |
| Monitoramento    | Prometheus + Grafana           |
| CI/CD            | GitHub Actions                 |
| Orquestração     | Docker Compose (dev) / Kubernetes (prod) |

### Fluxo resumido:

1. Usuário se autentica e envia um vídeo.
2. Sistema registra a requisição no PostgreSQL e publica um evento no RabbitMQ.
3. Microsserviço de processamento consome a fila, extrai os frames e gera o `.zip`.
4. Zip é armazenado (ex: S3 ou volume persistente).
5. Status do vídeo é atualizado no banco.
6. Usuário pode listar seus vídeos e baixar o ZIP.
7. Em caso de erro, um e-mail é disparado.

---

## ✅ Requisitos Atendidos

### Funcionais

- [x] Processamento simultâneo (vários consumidores na fila)
- [x] Mensageria garante que nenhuma requisição seja perdida
- [x] Autenticação por usuário/senha
- [x] Listagem de status por usuário
- [x] Notificação por e-mail em caso de erro

### Técnicos

- [x] Persistência de dados (PostgreSQL)
- [x] Arquitetura escalável horizontalmente
- [x] Versionamento no GitHub
- [x] Testes automatizados (unitários e integração)
- [x] CI/CD com GitHub Actions

---

## 🧪 Stack Tecnológica

| Categoria          | Escolha                          |
|--------------------|----------------------------------|
| Backend (upload)   | Java Spring Boot / Node.js       |
| Processamento      | Python + OpenCV                  |
| Mensageria         | RabbitMQ                         |
| Banco de dados     | PostgreSQL                       |
| Cache              | Redis                            |
| Containerização    | Docker                           |
| Orquestração       | Kubernetes (ou Docker Compose)   |
| Monitoramento      | Prometheus + Grafana             |
| CI/CD              | GitHub Actions                   |
| Email              | SMTP (SendGrid / Mailtrap)       |
