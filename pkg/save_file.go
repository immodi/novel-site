package pkg

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile saves an uploaded file to the target directory and returns its relative path
func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, targetDir string) (string, error) {
	// Ensure directory exists
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", err
	}

	// Create unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(targetDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	// Return relative path (to be stored in DB)
	return "/" + filePath, nil
}
