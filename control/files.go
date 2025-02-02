package control

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/models"
	"sifu-box/utils"

	"go.uber.org/zap"
)

func GetFiles(privateKey, workDir string, entClient *ent.Client, logger *zap.Logger) (map[string][]map[string]string, []error){
	var errors []error
	fileMap := make(map[string]string)
	token, err := utils.EncryptionMd5(privateKey)
	if err != nil {
		logger.Error(fmt.Sprintf("计算md5密钥失败: [%s]", err.Error()))
		errors = append(errors, fmt.Errorf("计算md5密钥失败"))
		return nil, errors
	}
	providers, err := entClient.Provider.Query().Select(provider.FieldName).All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("查询数据库数据失败: [%s]", err.Error()))
		errors = append(errors, fmt.Errorf("查询数据库机场数据失败"))
		return nil, errors
	}
	for _, provider := range providers {
		providerHashName, err := utils.EncryptionMd5(provider.Name)
		if err != nil {
			logger.Error(fmt.Sprintf("计算'%s'机场的md5失败: [%s]", provider.Name, err.Error()))
			errors = append(errors, fmt.Errorf("计算'%s'机场的md5失败", provider.Name))
			continue
		}
		fileMap[fmt.Sprintf("%s.json", providerHashName)] = provider.Name
	}
	dirs ,err := os.ReadDir(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR))
	if err != nil {
		logger.Error(fmt.Sprintf("遍历配置文件夹失败: [%s]", err.Error()))
		errors = append(errors, fmt.Errorf("遍历配置文件夹失败"))
		return nil, errors
	}
	fileGroup := make(map[string][]map[string]string)
	for _, dir := range dirs {
		if !dir.IsDir() {
			logger.Error(fmt.Sprintf("配置文件夹下的模板'%s'不是文件夹", dir.Name()))
			errors = append(errors, fmt.Errorf("配置文件夹下的模板'%s'不是文件夹", dir.Name()))
			continue
		}
		
		files, err := os.ReadDir(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR, dir.Name()))
		if err != nil {
			logger.Error(fmt.Sprintf("遍历'%s'模板文件夹失败: [%s]", dir.Name(), err.Error()))
			errors = append(errors, fmt.Errorf("遍历'%s'模板文件夹失败", dir.Name()))
			continue
		}
		var fileLinks []map[string]string
		for _, file := range files {
			name, ok := fileMap[file.Name()]
			if !ok {
				logger.Error(fmt.Sprintf("文件'%s'对应的机场不存在", file.Name()))
				errors = append(errors, fmt.Errorf("文件'%s'对应的机场不存在", file.Name()))
				continue
			}
			path, err := url.JoinPath("api", "file", name)
			if err != nil {
				logger.Error(fmt.Sprintf("拼接'%s'文件路径失败: [%s]", name, err.Error()))
				errors = append(errors, fmt.Errorf("拼接'%s'文件路径失败", name))
				continue
			}
			params := url.Values{}
			params.Add("token", token)
			params.Add("template", dir.Name())
			params.Add("path", file.Name())
			path += "?" + params.Encode()
			fileLinks = append(fileLinks, map[string]string{"label": name, "path": path})
		}
		fileGroup[dir.Name()] = fileLinks
	}
	return fileGroup, nil
}