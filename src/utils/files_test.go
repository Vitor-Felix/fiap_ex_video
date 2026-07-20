package utils

import (
	"os"
	"testing"
)

func TestIsValidVideoFile(t *testing.T) {
	// Padrão Table-Driven Tests: definimos um slice de structs com os cenários
	cenarios := []struct {
		nome     string
		arquivo  string
		esperado bool
	}{
		{"Extensão válida minúscula", "meu_video.mp4", true},
		{"Extensão válida maiúscula", "AULA.AVI", true},
		{"Extensão inválida imagem", "foto.png", false},
		{"Arquivo sem extensão", "video_sem_nada", false},
		{"Extensão válida webm", "filme.webm", true},
	}

	for _, c := range cenarios {
		t.Run(c.nome, func(t *testing.T) {
			resultado := IsValidVideoFile(c.arquivo)
			if resultado != c.esperado {
				// No Go, usamos Errorf para formatar a string de falha e falhar o teste
				t.Errorf("Para o arquivo '%s': esperado %v, recebido %v", c.arquivo, c.esperado, resultado)
			}
		})
	}
}

func TestCreateDirs(t *testing.T) {
	// Muda o BasePath temporariamente para o teste não sujar a raiz do seu projeto
	BasePath = "./test_dir/"
	defer os.RemoveAll(BasePath) // Limpa tudo quando o teste terminar

	// Executa a função
	CreateDirs()

	// Verifica se as pastas foram realmente criadas
	pastasEsperadas := []string{
		BasePath + "uploads",
		BasePath + "outputs",
		BasePath + "temp",
	}

	for _, pasta := range pastasEsperadas {
		if _, err := os.Stat(pasta); os.IsNotExist(err) {
			t.Errorf("A pasta %s não foi criada", pasta)
		}
	}
}
