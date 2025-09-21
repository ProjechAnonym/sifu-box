package control

import (
	"context"
	"fmt"
	"os"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/ruleset"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FetchItems(ent_client *ent.Client, logger *zap.Logger) (gin.H, error) {
	providers, err := ent_client.Provider.Query().Select(provider.FieldID, provider.FieldName, provider.FieldPath, provider.FieldRemote).All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取机场列表失败: [%s]", err.Error()))
		return gin.H{"message": fmt.Sprintf("获取机场列表失败: [%s]", err.Error())}, err
	}
	rulesets, err := ent_client.Ruleset.Query().Select(ruleset.FieldID, ruleset.FieldName, ruleset.FieldPath, ruleset.FieldRemote, ruleset.FieldUpdateInterval, ruleset.FieldBinary, ruleset.FieldDownloadDetour).All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取规则集列表失败: [%s]", err.Error()))
		return gin.H{"message": fmt.Sprintf("获取规则集列表失败: [%s]", err.Error())}, err
	}
	items := struct {
		Providers []struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Path   string `json:"path"`
			Remote bool   `json:"remote"`
		} `json:"providers"`
		Rulesets []struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Path           string `json:"path"`
			Remote         bool   `json:"remote"`
			UpdateInterval string `json:"update_interval"`
			Binary         bool   `json:"binary"`
			DownloadDetour string `json:"download_detour"`
		} `json:"rulesets"`
	}{}
	for _, provider := range providers {
		items.Providers = append(items.Providers, struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Path   string `json:"path"`
			Remote bool   `json:"remote"`
		}{
			ID:     provider.ID,
			Name:   provider.Name,
			Path:   provider.Path,
			Remote: provider.Remote,
		})
	}
	for _, ruleset := range rulesets {
		items.Rulesets = append(items.Rulesets, struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Path           string `json:"path"`
			Remote         bool   `json:"remote"`
			UpdateInterval string `json:"update_interval"`
			Binary         bool   `json:"binary"`
			DownloadDetour string `json:"download_detour"`
		}{
			ID:             ruleset.ID,
			Name:           ruleset.Name,
			Path:           ruleset.Path,
			Remote:         ruleset.Remote,
			UpdateInterval: ruleset.UpdateInterval,
			Binary:         ruleset.Binary,
			DownloadDetour: ruleset.DownloadDetour,
		})
	}
	return gin.H{"message": items}, nil
}
func AddProvider(providers []model.Provider, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	for _, provider := range providers {
		if err := ent_client.Provider.Create().SetName(provider.Name).SetPath(provider.Path).SetRemote(provider.Remote).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加机场"%s"失败: [%s]`, provider.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加机场"%s"失败: [%s]`, provider.Name, err.Error())})
		} else {
			res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加机场"%s"成功`, provider.Name)})
		}
	}
	return res
}
func AddRuleset(rulesets []model.Ruleset, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	for _, ruleset := range rulesets {
		if err := ent_client.Ruleset.Create().SetName(ruleset.Name).SetPath(ruleset.Path).SetRemote(ruleset.Remote).SetBinary(ruleset.Binary).SetDownloadDetour(ruleset.DownloadDetour).SetUpdateInterval(ruleset.UpdateInterval).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加规则集"%s"失败: [%s]`, ruleset.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加规则集"%s"失败: [%s]`, ruleset.Name, err.Error())})
		} else {
			res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加规则集"%s"成功`, ruleset.Name)})
		}
	}
	return res
}
func EditProvider(name, path string, remote bool, ent_client *ent.Client, logger *zap.Logger) error {

	if err := ent_client.Provider.Update().Where(provider.NameEQ(name)).SetPath(path).SetRemote(remote).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf(`修改机场"%s"失败: [%s]`, name, err.Error()))
	}
	return nil
}
func EditRuleset(name, path, update_interval, download_detour string, remote, binary bool, ent_client *ent.Client, logger *zap.Logger) error {
	if err := ent_client.Ruleset.Update().Where(ruleset.NameEQ(name)).SetPath(path).SetDownloadDetour(download_detour).SetUpdateInterval(update_interval).SetBinary(binary).SetRemote(remote).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf(`修改机场"%s"失败: [%s]`, name, err.Error()))
	}
	return nil
}
func DeleteProvider(name []string, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	for _, n := range name {
		provider_msg, err := ent_client.Provider.Query().Where(provider.NameEQ(n)).First(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf(`查找机场"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找机场"%s"失败: [%s]`, n, err.Error())})
			continue
		}
		if !provider_msg.Remote {
			if err := os.Remove(provider_msg.Path); err != nil {
				logger.Error(fmt.Sprintf(`删除机场"%s"文件失败: [%s]`, n, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除机场"%s"文件失败: [%s]`, n, err.Error())})
				continue
			}
		}
		if _, err := ent_client.Provider.Delete().Where(provider.NameEQ(n)).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`删除机场"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除机场"%s"失败: [%s]`, n, err.Error())})
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`删除机场"%s"成功`, n)})
	}
	return res
}
func DeleteRuleset(name []string, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	for _, n := range name {
		rule_set, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(n)).First(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf(`查找规则集"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找规则集"%s"失败: [%s]`, n, err.Error())})
			continue
		}
		if !rule_set.Remote {
			if err := os.Remove(rule_set.Path); err != nil {
				logger.Error(fmt.Sprintf(`删除规则集"%s"文件失败: [%s]`, n, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除规则集"%s"文件失败: [%s]`, n, err.Error())})
				continue
			}
		}
		if _, err := ent_client.Ruleset.Delete().Where(ruleset.NameEQ(n)).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`删除规则集"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除规则集"%s"失败: [%s]`, n, err.Error())})
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`删除规则集"%s"成功`, n)})
	}
	return res
}
