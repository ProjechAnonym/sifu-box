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
	UpdateInterval string `json:"update_interval,omitempty" yaml:"update_interval,omitempty"`
	Binary         bool   `json:"binary" yaml:"binary"`
	DownloadDetour string `json:"download_detour,omitempty" yaml:"download_detour,omitempty"`
}

// AutoFill 根据上传的规则集文件自动填充Ruleset结构体的字段
// 参数:
//
//	file: 上传的文件头信息, 包含文件名等元数据
//	work_dir: 工作目录路径, 用于构建文件保存路径
//
// 返回值:
//
//	error: 文件处理过程中的错误信息, 如果文件格式不支持则返回相应错误
func (r *Ruleset) AutoFill(file *multipart.FileHeader, work_dir string) error {
	// 提取文件扩展名和文件名
	ext := filepath.Ext(file.Filename)
	file_name := strings.TrimSuffix(file.Filename, ext)

	// 构建文件保存路径, 使用文件名的MD5值作为实际存储文件名
	save_path := filepath.Join(work_dir, "temp", "rulesets", fmt.Sprintf("%x%s", md5.Sum([]byte(file_name)), ext))

	// 根据文件扩展名判断规则集类型并设置二进制标志
	switch ext {
	case ".srs":
		r.Binary = true
	case ".json":
		r.Binary = false
	default:
		return fmt.Errorf(`"%s"不支持的规则集格式: [%s]`, file_name, ext)
	}

	// 设置规则集的基本信息
	r.Name = file_name
	r.Path = save_path
	r.Remote = false
	return nil
}
