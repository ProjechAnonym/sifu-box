package utils

import (
	"io/fs"
	"os"
	"path/filepath"
)

func File_write(content []byte, dst string,perm []fs.FileMode) error {
	if _, err := os.Stat(filepath.Dir(dst)); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(dst), perm[0])
		}
	}
	_, err := os.Stat(dst)
	if err != nil {
		if os.IsNotExist(err) {
			os.Remove(dst)
		}
	} else {
		os.Remove(dst)
	}
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, perm[1])
	defer func() {
		if err := file.Close(); err != nil {
			Logger_caller("File can not close!", err,1)
		}
	}()
	if err != nil {
		Logger_caller("Create file failed", err,1)
		return err
	}
	_, err = file.WriteString(string(content))
	if err != nil {
		Logger_caller("Write config failed!", err,1)
		return err
	}
	return nil
}
