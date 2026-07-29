# ☸️ Guia de Execução Local (Kubernetes via Minikube)

Este documento descreve o passo a passo para executar o sistema FIAP X utilizando Kubernetes via Minikube.

A arquitetura executada no cluster contém:

- PostgreSQL
- RabbitMQ
- API Backend Go
- Worker Python consumidor de fila
- Nginx Gateway
- Prometheus para coleta de métricas
- Grafana para visualização

## Arquitetura do Sistema

```text
                +------------------+
                |     Gateway      |
                +---------+--------+
                          |
                          v
                   +-------------+
                   | API (Go)    |
                   +------+------+ 
                          |
          +---------------+---------------+
          |                               |
          v                               v
   PostgreSQL                      RabbitMQ
                                          |
                                          v
                                   Worker (Python)

                     Prometheus ---> API /metrics
                     Prometheus ---> Worker /metrics

                     Grafana -----> Prometheus
```

---

## 📋 Pré-requisitos

Certifique-se de possuir:

- Docker instalado e em execução
- kubectl instalado
- Minikube instalado

Valide com:

```bash
docker --version
kubectl version --client
minikube version
```

---

## 🚀 1. Inicializar o Cluster Kubernetes

Inicie o Minikube:

```bash
minikube start
```

Valide:

```bash
kubectl get nodes
```

Resultado esperado:

```text
NAME       STATUS   ROLES
minikube   Ready    control-plane
```

---

## 🐳 2. Construir as Imagens Docker

O Minikube possui um ambiente Docker próprio. Portanto, as imagens precisam ser carregadas dentro do cluster.

> ⚠️ Durante o desenvolvimento, recomenda-se utilizar versões incrementais das imagens (v1, v2, v3...) em vez da tag latest. Isso evita que o Minikube reutilize imagens antigas e dificulte a validação de alterações recentes.

### API Go

```bash
docker build -t fiap-api:v3 -f Dockerfile .
```

### Worker Python

```bash
docker build -t fiap-worker:v4 -f worker/Dockerfile ./worker
```

Valide localmente:

```bash
docker images | grep fiap
```

Esperado:

```text
fiap-api       v3
fiap-worker    v4
```

---

## 📦 3. Carregar imagens no Minikube

Envie as imagens para dentro do cluster:

```bash
minikube image load fiap-api:v3
minikube image load fiap-worker:v4
```

Valide:

```bash
minikube image ls | grep fiap
```

Esperado:

```text
fiap-api:v3
fiap-worker:v4
```

---

## ⚙️ 4. Configurações auxiliares

Crie os ConfigMaps utilizados pelo Kubernetes.

### Nginx

```bash
kubectl create configmap nginx-config \
  --from-file=nginx/nginx.conf
```

### PostgreSQL

```bash
kubectl create configmap postgres-init-script \
  --from-file=db/init.sql
```

Se eles já existirem:

```bash
kubectl delete configmap nginx-config postgres-init-script
```

e execute novamente os comandos acima.

---

## ☸️ 5. Aplicar os manifestos Kubernetes

Aplique todos os recursos:

```bash
kubectl apply -f k8s/
```

> ⚠️ **Credenciais Mailtrap (notificação por e-mail):** o arquivo `k8s/00-configmap-secrets.yaml` contém
> placeholders `CHANGE_ME` para as variáveis SMTP. Após o `kubectl apply -f k8s/`, atualize o Secret
> com as credenciais reais usando o comando abaixo — ele não altera nenhum arquivo do repositório:
>
> ```bash
> kubectl patch secret fiap-secrets \
>   --type='json' \
>   -p='[
>     {"op":"replace","path":"/data/SMTP_USER","value":"'$(echo -n "SEU_SMTP_USER" | base64)'"},
>     {"op":"replace","path":"/data/SMTP_PASSWORD","value":"'$(echo -n "SUA_SMTP_PASSWORD" | base64)'"}
>   ]'
> ```
>
> Depois reinicie o worker para ele carregar os novos valores:
>
> ```bash
> kubectl rollout restart deployment/fiap-worker-deployment
> ```

Acompanhe a inicialização:

```bash
kubectl get pods
```

Todos os pods devem ficar com status `Running` e `READY 1/1`.

---

## 🔎 6. Validar comunicação interna

Verifique os services:

```bash
kubectl get svc
```

Verifique os endpoints:

```bash
kubectl get endpoints
```

---

## 🌐 7. Acessar a aplicação

Crie um tunnel para o gateway:

```bash
kubectl port-forward svc/gateway-service 8080:8080
```

Depois abra no navegador:

```text
http://localhost:8080
```

---

## 🗄️ 8. Acessar o PostgreSQL

Em outro terminal:

```bash
kubectl port-forward svc/postgres-service 5432:5432
```

Conexão com DBeaver ou psql:

```text
Host: localhost
Port: 5432
Database: fiap_x_db
User: fiap_user
```

---

## 🗃️ 8.5. Migration — adicionar coluna email

> ⚠️ Este passo é necessário se o volume do PostgreSQL já existia antes desta versão.
> Se o cluster foi criado do zero com o `init.sql` atual, a coluna já existe e pode pular.

```bash
kubectl exec deployment/postgres-deployment -- psql -U fiap_user -d fiap_x_db \
  -c "ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255);"
```

---

## 👤 9. Criar usuários iniciais

Como ainda não existe tela de cadastro, execute o seguinte SQL no banco:

```sql
INSERT INTO users (id, username, password_hash, email, created_at)
VALUES 
    (gen_random_uuid(), 'admin', crypt('123456', gen_salt('bf')), 'admin@fiap-x.local', NOW()),
    (gen_random_uuid(), 'dev',   crypt('123456', gen_salt('bf')), 'dev@fiap-x.local',   NOW());
```

> 💡 **Para testar notificações por e-mail:** substitua o campo `email` pelo endereço real da sua
> conta Mailtrap. Qualquer falha de processamento disparará o alerta para esse endereço:
>
> ```sql
> UPDATE users SET email = 'seu@email.com' WHERE username = 'admin';
> ```

Credenciais:

```text
admin / 123456
dev   / 123456
```

---

## 📊 10. Prometheus

A aplicação expõe métricas em:

- API: `/metrics`
- Worker: `:8000/metrics`

Abra o Prometheus:

```bash
kubectl port-forward svc/prometheus-service 9090:9090
```

Acesse:

```text
http://localhost:9090
```

Na interface, valide o painel de `Targets` e confirme se os alvos aparecem como `UP`.

### Validação direta das métricas no cluster

Também é possível validar diretamente do cluster para identificar se o problema está no Prometheus, no Service, no Endpoint ou na aplicação:

```bash
kubectl run curl-test \
  --image=curlimages/curl \
  -it --rm -- sh
```

Dentro do container:

```bash
curl http://api-service:8080/metrics
curl http://worker-service:8000/metrics
```

---

## 📈 11. Grafana

Abra o Grafana:

```bash
kubectl port-forward svc/grafana-service 3000:3000
```

Acesse:

```text
http://localhost:3000
```

Credenciais padrão:

```text
admin / admin
```

### Configuração inicial

1. Adicione uma Data Source do tipo Prometheus.
2. Informe a URL:

```text
http://prometheus-service:9090
```

3. Salve.

O dashboard do projeto pode ser importado a partir do arquivo JSON disponível no repositório. Acesse `Dashboards → Import → Upload JSON` e selecione o arquivo.

---

## 🛠️ Troubleshooting

### Problema: Alterei o código, mas o Kubernetes continua usando a versão antiga

Evite usar a tag `latest` durante o desenvolvimento. Prefira versões incrementais como `v1`, `v2`, `v3`.

Depois de gerar uma nova imagem:

```bash
docker build -t fiap-api:v3 .
minikube image load fiap-api:v3
```

Atualize a tag no Deployment e aplique:

```bash
kubectl apply -f k8s/
kubectl rollout restart deployment fiap-api-deployment
```

### Problema: Prometheus encontra o serviço, mas recebe HTTP 404

**Sintoma:** o target fica `DOWN` com erro `404 Not Found`.

**Causa:** a API em execução não possui o endpoint `/metrics`, geralmente porque uma imagem Docker antiga foi carregada no cluster.

**Validação:**

```bash
kubectl logs deploy/fiap-api-deployment
```

Se aparecer algo como `GET /metrics 404`, a imagem está desatualizada.

**Solução:** gere uma nova tag da imagem, atualize o Deployment e reinicie o rollout.

### Problema: ConfigMap alterado, mas o Prometheus continua usando a configuração antiga

Reinicie o deployment:

```bash
kubectl rollout restart deployment prometheus-deployment
```

### Problema: Porta já utilizada no port-forward

Exemplo de erro:

```text
address already in use
```

Use outra porta:

```bash
kubectl port-forward svc/prometheus-service 9091:9090
```

---

## 🧹 Encerrar o ambiente Kubernetes

Para parar o cluster:

```bash
minikube stop
```

Para remover completamente:

```bash
minikube delete
```

