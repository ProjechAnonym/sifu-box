package utils

import (
	"os"
	"path/filepath"
)

func WriteFile(path string, data []byte, tag int, perm os.FileMode) error {
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), 0755)
	}
	file, err := os.OpenFile(path, tag, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)
	return nil
}