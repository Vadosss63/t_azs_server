package application

import (
	"io"
	"os"
	"path/filepath"
)

func EnsureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func ListFilesInDirectory(directory string) ([]string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil
}

func DeleteDirectory(directoryPath string) error {
	return os.RemoveAll(directoryPath)
}

func SaveUploadedFile(uploadsDir, filename string, file io.Reader) error {
	if err := EnsureDirectory(uploadsDir); err != nil {
		return err
	}

	safeFilename := filepath.Base(filename)
	dst, err := os.Create(filepath.Join(uploadsDir, safeFilename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return err
}

func CheckFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
