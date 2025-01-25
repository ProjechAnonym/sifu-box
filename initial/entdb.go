package initial

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/ruleset"
	"sifu-box/ent/template"
	"sifu-box/models"

	"entgo.io/ent/dialect"
	"go.uber.org/zap"
)

func InitEntdb(workDir string, logger *zap.Logger) (*ent.Client){
	entClient, err := ent.Open(dialect.SQLite, fmt.Sprintf("file:%s/sifu-box.db?cache=shared&_fk=1", workDir))
	if err != nil {
		logger.Error(fmt.Sprintf("连接Ent数据库失败: [%s]",err.Error()))
		panic(err)
	}
	logger.Info("连接Ent数据库完成")
	if err = entClient.Schema.Create(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("创建表资源失败: [%s]",err.Error()))
		panic(err)
	}
	logger.Info("自动迁移Ent数据库完成")
	return entClient
}

func SaveNewProxySetting(configuration models.Configuration, entClient *ent.Client, logger *zap.Logger) {
	for _, supplier := range configuration.Providers {
		exist, err := entClient.Provider.Query().Where(provider.NameEQ(supplier.Name)).Exist(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		}
		if !exist {
			if _, err := entClient.Provider.Create().SetName(supplier.Name).SetDetour(supplier.Detour).SetPath(supplier.Path).SetRemote(supplier.Remote).Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
			}
		}	
	}
	logger.Info("数据库写入机场信息完成")

	for _, collectionInfo := range configuration.Rulesets {
		exist, err := entClient.RuleSet.Query().Where(ruleset.TagEQ(collectionInfo.Tag)).Exist(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		}
		if !exist {
			if _, err := entClient.RuleSet.Create().SetTag(collectionInfo.Tag).
													SetNameServer(collectionInfo.NameServer).
													SetPath(collectionInfo.Path).
													SetType(collectionInfo.Type).
													SetFormat(collectionInfo.Format).
													SetChina(collectionInfo.China).
													SetLabel(collectionInfo.Label).
													SetDownloadDetour(collectionInfo.DownloadDetour).
													SetUpdateInterval(collectionInfo.UpdateInterval).
													Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
			}
		}
	}
	logger.Info("数据库写入规则集信息完成")

	for key, templateContent := range configuration.Templates {
		exist, err := entClient.Template.Query().Where(template.NameEQ(key)).Exist(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		}
		if !exist {
			if _, err := entClient.Template.Create().
											SetName(key).
											SetContent(templateContent).
											Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
			}
		}
	}
	logger.Info("数据库写入模板信息完成")
}