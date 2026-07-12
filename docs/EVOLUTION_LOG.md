# 📈 LOG DE EVOLUÇÃO E APRENDIZADO

Este arquivo documenta as alterações de engenharia feitas no projeto e os conceitos consolidados pelo desenvolvedor até o momento.

---

## 📌 [CHAT 1] - Finalização do Milestone 1 (Issues 1.1 e 1.2)

### 🧠 Conceitos Consolidados pelo Desenvolvedor:
1. **Gargalo Síncrono no Go:** Entendimento de como o monólito original bloqueia requisições HTTP ao chamar o `ffmpeg` via CLI (`os/exec`) de forma síncrona.
2. **Filosofia do Go:** Tratamento explícito de erros (`if err != nil`), uso do `defer` para fechamento de ponteiros de I/O, sintaxe de `structs` e Struct Tags para parsing JSON.
3. **Máquina de Estados:** A necessidade de usar estados (`PENDENTE`, `PROCESSANDO`, `CONCLUIDO`, `ERRO`) via ENUM do Postgres para rastrear processos assíncronos.
4. **Persistência via Docker Compose:** Configuração de infraestrutura reprodutível e isolada usando volumes locais (`pgdata`) para evitar a perda de dados ao reiniciar containers.

### 🛠️ Código Evoluído no Repositório:
* Criado o script `db/init.sql` com a modelagem da tabela `videos` utilizando UUID e ENUM de status.
* Criado o arquivo `docker-compose.yml` na raiz gerenciando o container do PostgreSQL e seu respectivo volume persistente.
* Testado e validado o funcionamento do banco local via CLI `psql`.

---
*(A próxima IA deverá continuar a partir deste ponto, iniciando a Issue 2.1)*
