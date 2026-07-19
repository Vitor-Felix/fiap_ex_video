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
