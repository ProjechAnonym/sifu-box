package model

import (
	"crypto/md5"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type Provider struct {
	Name   string `json:"name" yaml:"name"`
	Path   string `json:"path" yaml:"path"`
	Remote bool   `json:"remote" yaml:"remote"`
}

// AutoFill 根据上传的文件信息自动填充Provider结构体的字段
// 参数:
//
//	file: 上传的文件头信息, 包含文件名等元数据
//	work_dir: 工作目录路径, 用于构建文件保存路径
//
// 返回值:
//
//	error: 如果文件格式不支持则返回错误, 否则返回nil
func (p *Provider) AutoFill(file *multipart.FileHeader, work_dir string) error {
	// 提取文件扩展名和文件名
	ext := filepath.Ext(file.Filename)
	file_name := strings.TrimSuffix(file.Filename, ext)

	// 构建文件保存路径, 将文件名进行MD5加密后保存
	save_path := filepath.Join(work_dir, "temp", "providers", fmt.Sprintf("%x%s", md5.Sum([]byte(file_name)), ext))

	// 验证文件格式是否支持, 只允许yaml和yml格式
	switch ext {
	case ".yaml", ".yml":
	default:
		return fmt.Errorf("不支持的文件格式: [%s]", ext)
	}

	// 填充Provider结构体字段
	p.Name = file_name
	p.Path = save_path
	p.Remote = false
	return nil
}
