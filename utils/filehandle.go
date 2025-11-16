package utils

import (
	"io"
	"os"
	"path/filepath"
)

func WriteFile(path string, data []byte, tag int, perm os.FileMode) error {
	// 检查文件目录是否存在, 如果不存在则创建
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), 0755)
	}
	// 打开或创建文件
	file, err := os.OpenFile(path, tag, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入数据
	file.Write(data)
	return nil
}

func ReadFile(path string) ([]byte, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
