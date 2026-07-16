package utils

import (
	"os"
	"path/filepath"
)

// BasePath guarda o prefixo do caminho ("" se rodando da raiz, "../" se rodando de /src)
var BasePath = ""

func init() {
	// Buscamos uma variável de ambiente chamada APP_ENV.
	// Se ela for "production", usamos "" (diretório raiz do container).
	// Se não estiver setada, fazemos a checagem local baseada em onde o binário roda.
	if os.Getenv("APP_ENV") == "production" {
		BasePath = ""
		return
	}

	// Fallback local: se a pasta "handlers" não existir no diretório atual,
	// significa que estamos rodando de fora da pasta "src".
	if _, err := os.Stat("handlers"); os.IsNotExist(err) {
		BasePath = "../"
	} else {
		BasePath = ""
	}
}

// CreateDirs garante que os diretórios necessários existam
func CreateDirs() {
	dirs := []string{
		filepath.Join(BasePath, "uploads"),
		filepath.Join(BasePath, "outputs"),
		filepath.Join(BasePath, "temp"),
	}
	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}
}
