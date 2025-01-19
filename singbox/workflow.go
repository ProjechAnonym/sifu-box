package singbox

import (
	"context"
	"encoding/json"
	"fmt"
	"sifu-box/ent"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func Workflow(entClient *ent.Client, buntClient *buntdb.DB, logger *zap.Logger) ([]string, error) {
	settingStr, err := utils.GetValue(buntClient, "setting", logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取设置信息失败: [%s]", err.Error()))
		return nil, err
	}
	var setting models.Setting
	if err := json.Unmarshal([]byte(settingStr), &setting); err != nil {
		logger.Error(fmt.Sprintf("解析设置信息失败: [%s]", err.Error()))
		return nil, err
	}
	var providers []models.Provider
	if setting.Server.Enabled {
		providerList, err := entClient.Provider.Query().All(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取机场信息失败: [%s]", err.Error()))
			return nil, err
		}
		for _, provider := range providerList {
			providers = append(providers, models.Provider{
				Name: provider.Name,
				Path: provider.Path,
				Remote: provider.Remote,
				Detour: provider.Detour,
			})
		}
	}else{
		providers = setting.Providers
	}
	merge(providers, logger)
	// templateMap := setting.Templates

	

	return nil, nil
}