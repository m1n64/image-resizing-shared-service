package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

var mimeExtMap = map[string]string{
	"image/jpeg": ".jpg",
	"image/jpg":  ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

func GetAppDir() string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(exe)
}

func GetDataDirPath() string {
	return filepath.Join(GetAppDir(), "data")
}

func GetDataPath(file string) string {
	return filepath.Join(GetDataDirPath(), file)
}

func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func ReadFileFromStorage(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", path, err)
	}
	return data, nil
}

func SaveFileToStorage(filePath, destPath string) error {
	srcFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("cannot create dest directory: %w", err)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	return nil
}

func SaveBytesToStorage(filePath string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	}

	return nil
}

func SaveTempFile(file io.Reader) (string, string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", "", err
	}
	contentType := http.DetectContentType(buf[:n])

	ext := mimeExtMap[contentType]
	if ext == "" {
		if exts, err := mime.ExtensionsByType(contentType); err == nil && len(exts) > 0 {
			ext = exts[0]
		} else {
			ext = ".bin"
		}
	}

	reader := io.MultiReader(bytes.NewReader(buf[:n]), file)

	tempFile, err := os.CreateTemp("tmp", "upload-*"+ext)
	if err != nil {
		return "", "", err
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		return "", "", err
	}

	return tempFile.Name(), contentType, nil
}

func RemoveFile(path string) {
	_ = os.Remove(path)
}
