package control

import (
	"context"
	"crypto/md5"
	"fmt"
	"path/filepath"
	"sifu-box/ent"
	"sifu-box/utils"

	"go.uber.org/zap"
)

func FileList(ent_client *ent.Client, logger *zap.Logger) ([]map[string]string, error) {
	templates, err := ent_client.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("查询模板失败: [%s]", err.Error()))
		return nil, err
	}
	template_list := []map[string]string{}
	for _, v := range templates {
		path := fmt.Sprintf("%x.json", md5.Sum([]byte(v.Name)))
		template_list = append(template_list, map[string]string{"name": v.Name, "path": path})
	}
	return template_list, nil
}
func FileDownload(work_dir, path string, logger *zap.Logger) ([]byte, error) {
	file_path := filepath.Join(work_dir, "sing-box", "config", path)
	content, err := utils.ReadFile(file_path)
	if err != nil {
		logger.Error(fmt.Sprintf("读取文件失败: [%s]", err.Error()))
		return nil, err
	}
	return content, nil
}
