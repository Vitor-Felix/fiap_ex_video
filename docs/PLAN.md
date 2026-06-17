🏁 Milestone 1: Entendimento e Infraestrutura Base
Objetivo: Rodar o projeto original, entender o fluxo e preparar o ambiente de banco de dados.

[ ] Issue 1.1: Executar e testar o projeto original localmente

O que fazer: Rodar o Docker atual, subir um vídeo e ver o ZIP ser gerado.

Foco de aprendizado: Entender como o Go usa o ffmpeg e onde os arquivos estão sendo salvos.

[ ] Issue 1.2: Modelagem e Criação do Banco de Dados (PostgreSQL)

O que fazer: Criar um container Docker isolado com PostgreSQL e desenhar a tabela de controle de vídeos (ID, nome do arquivo, status como PENDENTE, PROCESSANDO, CONCLUIDO, ERRO).

Foco de aprendizado: Persistência de dados e estados de um processo assíncrono.

🔑 Milestone 2: Persistência e Segurança (Ainda no Go)
Objetivo: Adicionar as primeiras regras de negócio essenciais antes de quebrar o sistema em microsserviços.

[ ] Issue 2.1: Conectar o código Go ao PostgreSQL

O que fazer: Alterar o main.go para que, ao receber o upload, ele salve o registro no banco antes de começar a processar.

Foco de aprendizado: Manipulação de banco de dados SQL (mesmo sem conhecer Go profundamente, vamos focar na lógica).

[ ] Issue 2.2: Implementar Autenticação Básica

O que fazer: Criar uma tela/rota de login simples. Proteger as rotas de upload e listagem para que apenas usuários logados vejam seus próprios vídeos.

Foco de aprendizado: Conceitos de Sessão, JWT ou autenticação simples.

🔀 Milestone 3: Desacoplamento com Mensageria (O Coração do Desafio)
Objetivo: Separar o recebimento do vídeo do processamento real para que o sistema não perca requisições em momentos de pico.

[ ] Issue 3.1: Subir o RabbitMQ e criar a Fila

O que fazer: Adicionar o RabbitMQ ao seu ambiente Docker.

Foco de aprendizado: Como funciona uma fila (Producer/Consumer) e por que ela garante a resiliência.

[ ] Issue 3.2: Transformar a API Go em um Produtor de Mensagens

O que fazer: Mudar o código Go. Agora, ao receber o vídeo, ele não vai mais processar. Ele apenas salva no banco (Status: PENDENTE), joga o ID do vídeo na fila do RabbitMQ e responde "Vídeo recebido!" para o usuário.

Foco de aprendizado: Comunicação assíncrona.

🐍 Milestone 4: O Novo Worker de Processamento
Objetivo: Criar o segundo microsserviço (pode ser em Python, aproveitando a sugestão do resumo, que é excelente para manipulação de vídeo).

[ ] Issue 4.1: Criar o Worker Consumidor

O que fazer: Criar um script que fica escutando a fila do RabbitMQ. Quando chega uma mensagem, ele altera o status no banco para PROCESSANDO.

Foco de aprendizado: Criação de microsserviços especialistas.

[ ] Issue 4.2: Extração de Frames e Geração do ZIP

O que fazer: Fazer o Worker executar o processamento do vídeo, gerar o ZIP e atualizar o status no banco para CONCLUIDO.

Foco de aprendizado: Processamento em background e compartilhamento de arquivos entre serviços (volumes Docker).

🛠️ Milestone 5: Qualidade, Monitoramento e CI/CD
Objetivo: Deixar o projeto pronto para a "produção" exigida pelos investidores.

[ ] Issue 5.1: Notificação de Erro (E-mail)

O que fazer: Configurar um serviço simulado de e-mail (como Mailtrap) para disparar um alerta se o processamento em background falhar.

[ ] Issue 5.2: Otimização dos Dockerfiles e Docker Compose Final

O que fazer: Consertar os Dockerfiles (remover os anti-padrões do original) e unir tudo (API, Worker, Postgres, RabbitMQ) em um único comando docker compose up.

[ ] Issue 5.3: Automação com GitHub Actions

O que fazer: Criar uma esteira simples de CI que roda testes ou valida o build a cada commit.