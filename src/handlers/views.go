package handlers

// GetHTMLForm retorna o HTML da página inicial do FIAP X integrado com o banco de dados
func GetHTMLForm() string {
	return `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FIAP X - Processador de Vídeos</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            max-width: 800px; 
            margin: 50px auto; 
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 { 
            color: #333; 
            text-align: center;
            margin-bottom: 30px;
        }
        .upload-form {
            border: 2px dashed #ddd;
            padding: 30px;
            text-align: center;
            border-radius: 10px;
            margin: 20px 0;
        }
        input[type="file"] {
            margin: 20px 0;
            padding: 10px;
        }
        button {
            background: #007bff;
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover { background: #0056b3; }
        .result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 5px;
            display: none;
        }
        .success { background: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
        .loading { 
            text-align: center; 
            display: none;
            margin: 20px 0;
        }
        .files-list {
            margin-top: 30px;
        }
        .file-item {
            background: #f8f9fa;
            padding: 12px;
            margin: 8px 0;
            border-radius: 5px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            border: 1px solid #e9ecef;
        }
        .download-btn {
            background: #28a745;
            color: white;
            padding: 5px 15px;
            text-decoration: none;
            border-radius: 3px;
            font-size: 14px;
        }
        .download-btn:hover { background: #218838; }
        
        /* Badges de Status adicionadas para o Histórico Real */
        .status-badge {
            padding: 3px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: bold;
            margin-left: 10px;
        }
        .status-PROCESSANDO { background: #fff3cd; color: #856404; }
        .status-CONCLUIDO { background: #d4edda; color: #155724; }
        .status-ERRO { background: #f8d7da; color: #721c24; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎬 FIAP X - Processador de Vídeos</h1>
        <p style="text-align: center; color: #666;">
            Faça upload de um vídeo e receba um ZIP com todos os frames extraídos!
        </p>
        
        <form id="uploadForm" class="upload-form">
            <p><strong>Selecione um arquivo de vídeo:</strong></p>
            <input type="file" id="videoFile" accept="video/*" required>
            <br>
            <button type="submit">🚀 Processar Vídeo</button>
        </form>
        
        <div class="loading" id="loading">
            <p>⏳ Processando vídeo... Isso pode levar alguns minutos.</p>
        </div>
        
        <div class="result" id="result"></div>
        
        <div class="files-list">
            <h3>📊 Seu Histórico de Processamento (Banco de Dados):</h3>
            <div id="filesList">Carregando histórico...</div>
        </div>
    </div>

    <script>
        document.getElementById('uploadForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const fileInput = document.getElementById('videoFile');
            const file = fileInput.files[0];
            
            if (!file) {
                showResult('Selecione um arquivo de vídeo!', 'error');
                return;
            }
            
            const formData = new FormData();
            formData.append('video', file);
            
            showLoading(true);
            hideResult();
            
            try {
                const response = await fetch('/upload', {
                    method: 'POST',
                    body: formData
                });
                
                const result = await response.json();
                
                if (result.success) {
                    showResult(
                        result.message + 
                        '<br><br><a href="/download/' + result.zip_path + '" class="download-btn">⬇️ Download ZIP</a>',
                        'success'
                    );
                    loadFilesList();
                } else {
                    showResult('Erro: ' + result.message, 'error');
                }
            } catch (error) {
                showResult('Erro de conexão: ' + error.message, 'error');
            } finally {
                showLoading(false);
            }
        });
        
        function showResult(message, type) {
            const result = document.getElementById('result');
            result.innerHTML = message;
            result.className = 'result ' + type;
            result.style.display = 'block';
        }
        
        function hideResult() {
            document.getElementById('result').style.display = 'none';
        }
        
        function showLoading(show) {
            document.getElementById('loading').style.display = show ? 'block' : 'none';
        }
        
        // Nova função loadFilesList conectada à nossa API do Postgres
        async function loadFilesList() {
            try {
                const response = await fetch('/api/videos');
                const videos = await response.json();
                
                const filesList = document.getElementById('filesList');
                
                if (videos && videos.length > 0) {
                    filesList.innerHTML = videos.map(video => {
                        let actionHtml = '';
                        
                        // Decide o botão ou mensagem dependendo do status vindo do Postgres
                        if (video.status === 'CONCLUIDO') {
                            actionHtml = '<a href="/download/' + video.zip_path + '" class="download-btn">⬇️ Download ZIP</a>';
                        } else if (video.status === 'ERRO') {
                            actionHtml = '<span style="color: #721c24; font-size: 13px;" title="' + (video.error_message || '') + '">⚠️ Falhou (Passe o mouse)</span>';
                        } else {
                            actionHtml = '<span style="color: #666; font-size: 13px;">⏳ Processando...</span>';
                        }

                        return '<div class="file-item">' +
                            '<div>' +
                                '<strong>' + video.original_name + '</strong>' +
                                '<span class="status-badge status-' + video.status + '">' + video.status + '</span>' +
                            '</div>' +
                            '<div>' + actionHtml + '</div>' +
                            '</div>';
                    }).join('');
                } else {
                    filesList.innerHTML = '<p>Nenhum vídeo processado ainda.</p>';
                }
            } catch (error) {
                document.getElementById('filesList').innerHTML = '<p>Erro ao carregar arquivos do banco de dados.</p>';
            }
        }
        
        // Carregar lista de arquivos ao inicializar
        loadFilesList();

        // Configura Polling para atualizar o histórico a cada 5 segundos de forma silenciosa
        setInterval(loadFilesList, 5000);
    </script>
</body>
</html>`
}
