package initial

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sifu-box/ent"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func SetDefaultTemplate(workDir string, buntClient *buntdb.DB, logger *zap.Logger) error {
	file, err := os.Open(filepath.Join(workDir, models.STATICDIR, models.TEMPLATEDIR, models.DEFAULTTEMPLATEPATH))
	if err != nil {
		logger.Error(fmt.Sprintf("打开默认模板文件失败: [%s]",err.Error()))
		return err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Error(fmt.Sprintf("读取默认模板文件失败: [%s]",err.Error()))
		return err
	}
	var template models.Template
	if err := yaml.Unmarshal(content, &template); err != nil {
		logger.Error(fmt.Sprintf("解析默认模板文件失败: [%s]",err.Error()))
		return err
	}
	if err := utils.SetValue(buntClient, models.DEFAULTTEMPLATEKEY, string(content), logger); err != nil {
		logger.Error(fmt.Sprintf("默认模板文件写入buntDB失败: [%s]",err.Error()))
		return err
	}
	return nil
}

func InitSetting(confDir string, server bool, buntClient *buntdb.DB, logger *zap.Logger) (*models.Setting, error){
	file, err := os.Open(filepath.Join(confDir, models.SIFUBOXSETTINGFILE))
	if err != nil {
		logger.Error(fmt.Sprintf("打开默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Error(fmt.Sprintf("读取默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	var setting models.Setting
	if err := yaml.Unmarshal(content, &setting); err != nil {
		logger.Error(fmt.Sprintf("解析默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	if !server && setting.Configuration == nil {
		logger.Error("缺少'configuration'字段, 非服务器模式下必须包含'configuration'字段")
		panic(fmt.Errorf("缺少'configuration'字段, 非服务器模式下必须包含'configuration'字段"))
	}
	if setting.Configuration != nil {
		configurationByte, err := yaml.Marshal(setting.Configuration)
		if err != nil {
			logger.Error(fmt.Sprintf("序列化默认模板文件失败: [%s]",err.Error()))
			return nil, err
		}
		if err := utils.SetValue(buntClient, models.SINGBOXSETTINGKEY, string(configurationByte), logger); err != nil {
			logger.Error(fmt.Sprintf("默认模板文件写入buntDB失败: [%s]",err.Error()))
			return nil, err
		}
	}
	if server {
		applicationByte, err := yaml.Marshal(setting.Application)
		if err != nil {
			logger.Error(fmt.Sprintf("序列化默认模板文件失败: [%s]",err.Error()))
			return nil, err
		}
		if err := utils.SetValue(buntClient, models.SIFUBOXSETTINGKEY, string(applicationByte), logger); err != nil {
			logger.Error(fmt.Sprintf("默认模板文件写入buntDB失败: [%s]",err.Error()))
			return nil, err
		}
	}
	return &setting, nil
}

func SetDefautlApplication(entClient *ent.Client, buntClient *buntdb.DB, logger *zap.Logger) error {
	providers, err := entClient.Provider.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取机场数据失败: [%s]",err.Error()))
		return fmt.Errorf("获取机场数据失败")
	}
	if len(providers) != 0 {
		if err := utils.SetValue(buntClient, models.CURRENTPROVIDER, providers[0].Name, logger); err != nil {
			logger.Error(fmt.Sprintf("设置默认机场失败: [%s]",err.Error()))
			return fmt.Errorf("设置默认机场失败")
		}
	}
	
	templates, err := entClient.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取模板信息失败: [%s]", err.Error()))
		return fmt.Errorf("获取模板信息失败")
	}
	if len(templates) != 0 {
		if err := utils.SetValue(buntClient, models.CURRENTTEMPLATE, templates[0].Name, logger); err != nil {
			logger.Error(fmt.Sprintf("设置默认机场失败: [%s]",err.Error()))
			return fmt.Errorf("设置默认机场失败")
		}
	}
	return nil
}