document.addEventListener('DOMContentLoaded', () => {
    checkAuth();

    // Handler de Login
    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const usernameInput = document.getElementById('username').value;
        const passwordInput = document.getElementById('password').value;

        try {
            const response = await fetch('/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: usernameInput, password: passwordInput })
            });

            const data = await response.json();

            if (response.ok && data.token) {
                localStorage.setItem('jwt_token', data.token);
                localStorage.setItem('username', data.username || usernameInput);
                showResult('Login efetuado com sucesso!', 'success');
                checkAuth();
            } else {
                showResult(data.error || 'Falha ao realizar login', 'error');
            }
        } catch (err) {
            showResult('Erro de conexão com o servidor', 'error');
        }
    });

    // Handler de Logout
    document.getElementById('logoutBtn').addEventListener('click', () => {
        localStorage.removeItem('jwt_token');
        localStorage.removeItem('username');
        checkAuth();
    });

    // Handler de Upload de Vídeo
    document.getElementById('uploadForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const token = localStorage.getItem('jwt_token');
        const fileInput = document.getElementById('videoFile');
        const file = fileInput.files[0];

        if (!file) return;

        const formData = new FormData();
        formData.append('video', file);

        showLoading(true);
        hideResult();

        try {
            const response = await fetch('/upload', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            const result = await response.json();

            if (response.status === 401) {
                handleUnauthorized();
                return;
            }

            if (result.success) {
                showResult(
                    result.message + `<br><br><a href="#" onclick="downloadZip('${result.zip_path}')" class="download-btn">⬇️ Download ZIP</a>`,
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
});

function checkAuth() {
    const token = localStorage.getItem('jwt_token');
    const username = localStorage.getItem('username');

    if (token) {
        document.getElementById('loginSection').style.display = 'none';
        document.getElementById('appSection').style.display = 'block';
        document.getElementById('loggedUser').innerText = username;
        loadFilesList();
    } else {
        document.getElementById('loginSection').style.display = 'block';
        document.getElementById('appSection').style.display = 'none';
    }
}

async function loadFilesList() {
    const token = localStorage.getItem('jwt_token');
    if (!token) return;

    try {
        const response = await fetch('/api/videos', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (response.status === 401) {
            handleUnauthorized();
            return;
        }

        const videos = await response.json();
        const filesList = document.getElementById('filesList');

        if (videos && videos.length > 0) {
            filesList.innerHTML = videos.map(video => {
                let actionHtml = '';

                if (video.status === 'CONCLUIDO') {
                    actionHtml = `<button onclick="downloadZip('${video.zip_path}')" class="download-btn">⬇️ Download ZIP</button>`;
                } else if (video.status === 'ERRO') {
                    actionHtml = `<span style="color: #721c24; font-size: 13px;" title="${video.error_message || ''}">⚠️ Falhou</span>`;
                } else {
                    actionHtml = '<span style="color: #666; font-size: 13px;">⏳ Processando...</span>';
                }

                return `<div class="file-item">
                    <div>
                        <strong>${video.original_name}</strong>
                        <span class="status-badge status-${video.status}">${video.status}</span>
                    </div>
                    <div>${actionHtml}</div>
                </div>`;
            }).join('');
        } else {
            filesList.innerHTML = '<p>Nenhum vídeo processado ainda.</p>';
        }
    } catch (error) {
        document.getElementById('filesList').innerHTML = '<p>Erro ao carregar arquivos do banco de dados.</p>';
    }
}

// Faz o download autenticado via JWT
async function downloadZip(zipPath) {
    const token = localStorage.getItem('jwt_token');
    try {
        const response = await fetch(`/download/${zipPath}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.status === 401) {
            handleUnauthorized();
            return;
        }

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = zipPath;
        document.body.appendChild(a);
        a.click();
        a.remove();
    } catch (err) {
        alert('Erro ao realizar download do arquivo.');
    }
}

function handleUnauthorized() {
    alert('Sessão expirada. Faça login novamente.');
    localStorage.removeItem('jwt_token');
    localStorage.removeItem('username');
    checkAuth();
}

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

// Polling silencioso
setInterval(() => {
    if (localStorage.getItem('jwt_token')) {
        loadFilesList();
    }
}, 5000);
