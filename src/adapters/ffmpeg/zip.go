package ffmpeg

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// CreateZipFile agrupa uma lista de caminhos de arquivos em um único ZIP
func CreateZipFile(files []string, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	// Fallback para garantir que o arquivo seja fechado caso ocorra um return antecipado com erro
	defer func() { _ = zipFile.Close() }()

	zipWriter := zip.NewWriter(zipFile)
	// Fallback para garantir que o writer seja limpo em caso de erro
	defer func() { _ = zipWriter.Close() }()

	for _, file := range files {
		err := addFileToZip(zipWriter, file)
		if err != nil {
			return err
		}
	}

	// 1. Fecha o zipWriter explicitamente para garantir a escrita do cabeçalho e rodapé do arquivo ZIP
	if err := zipWriter.Close(); err != nil {
		return err
	}

	// 2. Flusha e fecha o arquivo no sistema de arquivos tratando erros de disco
	if err := zipFile.Close(); err != nil {
		return err
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	// Ignora o erro no defer de leitura pois o arquivo já foi lido por completo no final da função
	defer func() { _ = file.Close() }()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
