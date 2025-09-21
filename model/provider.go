package model

import (
	"crypto/md5"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type Provider struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Remote bool   `json:"remote"`
}

func (p *Provider) AutoFill(file *multipart.FileHeader, work_dir string) error {
	ext := filepath.Ext(file.Filename)
	file_name := strings.TrimSuffix(file.Filename, ext)
	save_path := filepath.Join(work_dir, "uploads", "providers", fmt.Sprintf("%x%s", md5.Sum([]byte(file_name)), ext))
	switch ext {
	case ".yaml", ".yml":
	default:
		return fmt.Errorf("不支持的文件格式: [%s]", ext)
	}
	p.Name = file_name
	p.Path = save_path
	p.Remote = false
	return nil
}
