package model

import (
	"crypto/md5"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type Ruleset struct {
	Name           string `json:"name" yaml:"name"`
	Path           string `json:"path" yaml:"path"`
	Remote         bool   `json:"remote" yaml:"remote"`
	UpdateInterval string `json:"update_interval" yaml:"update_interval"`
	Binary         bool   `json:"binary" yaml:"binary"`
	DownloadDetour string `json:"download_detour" yaml:"download_detour"`
}

func (r *Ruleset) AutoFill(file *multipart.FileHeader, work_dir string) error {
	ext := filepath.Ext(file.Filename)
	file_name := strings.TrimSuffix(file.Filename, ext)
	save_path := filepath.Join(work_dir, "uploads", "rulesets", fmt.Sprintf("%x%s", md5.Sum([]byte(file_name)), ext))
	switch ext {
	case ".srs":
		r.Binary = true
	case ".json":
		r.Binary = false
	default:
		return fmt.Errorf(`"%s"不支持的规则集格式: [%s]`, file_name, ext)
	}
	r.Name = file_name
	r.Path = save_path
	r.Remote = false
	return nil
}
