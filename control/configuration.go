package control

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/ruleset"
	"sifu-box/ent/template"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FetchItems(ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	providers, err := ent_client.Provider.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取机场列表失败: [%s]", err.Error()))
		res = append(res, gin.H{"status": false, "message": fmt.Sprintf("获取机场列表失败: [%s]", err.Error())})
		return res
	}
	provider_list := []struct {
		ID int `json:"id"`
		model.Provider
	}{}
	ruleset_list := []struct {
		ID int `json:"id"`
		model.Ruleset
	}{}
	template_list := []struct {
		ID int `json:"id"`
		model.Template
	}{}
	for _, provider := range providers {
		provider_list = append(provider_list, struct {
			ID int `json:"id"`
			model.Provider
		}{
			ID: provider.ID,
			Provider: model.Provider{
				Name:   provider.Name,
				Path:   provider.Path,
				Remote: provider.Remote,
			},
		})
	}
	res = append(res, gin.H{"status": true, "message": provider_list, "type": "provider"})
	rulesets, err := ent_client.Ruleset.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取规则集列表失败: [%s]", err.Error()))
		res = append(res, gin.H{"status": false, "message": fmt.Sprintf("获取规则集列表失败: [%s]", err.Error())})
		return res
	}
	for _, ruleset := range rulesets {
		ruleset_list = append(ruleset_list, struct {
			ID int `json:"id"`
			model.Ruleset
		}{
			ID: ruleset.ID,
			Ruleset: model.Ruleset{
				Name:           ruleset.Name,
				Path:           ruleset.Path,
				Remote:         ruleset.Remote,
				UpdateInterval: ruleset.UpdateInterval,
				Binary:         ruleset.Binary,
				DownloadDetour: ruleset.DownloadDetour,
			},
		})
	}
	res = append(res, gin.H{"status": true, "message": ruleset_list, "type": "ruleset"})
	templates, err := ent_client.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取模板列表失败: [%s]", err.Error()))
		res = append(res, gin.H{"status": false, "message": fmt.Sprintf("获取模板列表失败: [%s]", err.Error())})
		return res
	}
	for _, template := range templates {
		template_list = append(template_list, struct {
			ID int `json:"id"`
			model.Template
		}{
			ID: template.ID,
			Template: model.Template{
				Name:           template.Name,
				Ntp:            &template.Ntp,
				Inbounds:       template.Inbounds,
				OutboundsGroup: template.OutboundGroups,
				Providers:      template.Providers,
				DNS:            &template.DNS,
				Experiment:     &template.Experiment,
				Log:            &template.Log,
				Route:          &template.Route,
			},
		},
		)
	}
	res = append(res, gin.H{"status": true, "message": template_list, "type": "template"})
	return res
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
	exist, err := ent_client.Provider.Query().Where(provider.NameEQ(name)).Exist(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查找机场"%s"失败: [%s]`, name, err.Error()))
		return fmt.Errorf(`查找机场"%s"失败: [%s]`, name, err.Error())
	} else if !exist {
		logger.Error(fmt.Sprintf(`未找到机场"%s"`, name))
		return fmt.Errorf(`机场"%s"不存在`, name)
	}
	if err := ent_client.Provider.Update().Where(provider.NameEQ(name)).SetPath(path).SetRemote(remote).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf(`修改机场"%s"失败: [%s]`, name, err.Error()))
	}
	return nil
}
func EditRuleset(name, path, update_interval, download_detour string, remote, binary bool, ent_client *ent.Client, logger *zap.Logger) error {
	exist, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(name)).Exist(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查找规则集"%s"失败: [%s]`, name, err.Error()))
		return fmt.Errorf(`查找规则集"%s"失败: [%s]`, name, err.Error())
	} else if !exist {
		logger.Error(fmt.Sprintf(`未找到规则集"%s"`, name))
		return fmt.Errorf(`规则集"%s"不存在`, name)
	}
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
		rule_set, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(n)).Select(ruleset.FieldPath, ruleset.FieldRemote).First(context.Background())
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
func AddTemplate(template model.Template, ent_client *ent.Client, logger *zap.Logger) error {
	ent_template := ent_client.Template.Create()
	template.CreateFillFields(ent_template)
	if err := template.LinkProvidersTable(ent_client); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := template.LinkRulesetsTable(ent_client); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := template.EditProviders(); err != nil {
		return err
	}
	template.EditRulesets()
	if err := ent_template.SetUpdated(true).SetName(template.Name).SetProviders(template.Providers).SetInbounds(template.Inbounds).SetOutboundGroups(template.OutboundsGroup).Exec(context.Background()); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func DeleteTemplate(name []string, work_dir string, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	for _, n := range name {
		template_msg, err := ent_client.Template.Query().Where(template.NameEQ(n)).Select(template.FieldName, template.FieldProviders, template.FieldRoute).First(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf(`查找模板"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找模板"%s"失败: [%s]`, n, err.Error())})
			continue
		}
		template_instance := model.Template{Route: &template_msg.Route, Providers: template_msg.Providers, Name: template_msg.Name}
		if err := template_instance.UnLinkProvidersTable(ent_client); err != nil {
			logger.Error(err.Error())
			res = append(res, gin.H{"status": false, "message": err.Error()})
			continue
		}
		if err := template_instance.UnLinkRulesetsTable(ent_client); err != nil {
			logger.Error(err.Error())
			res = append(res, gin.H{"status": false, "message": err.Error()})
			continue
		}
		if _, err := os.Stat(filepath.Join(work_dir, "sing-box", "config", fmt.Sprintf(`%s.json`, fmt.Sprintf(`%x`, md5.Sum([]byte(template_msg.Name)))))); err == nil {
			if err := os.Remove(filepath.Join(work_dir, "sing-box", "config", fmt.Sprintf(`%s.json`, fmt.Sprintf(`%x`, md5.Sum([]byte(template_msg.Name)))))); err != nil {
				logger.Error(fmt.Sprintf(`删除模板"%s"配置文件失败: [%s]`, n, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除模板"%s"配置文件失败: [%s]`, n, err.Error())})
				continue
			}
		} else if !os.IsNotExist(err) {
			logger.Error(fmt.Sprintf(`查找模板"%s"配置文件失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找模板"%s"配置文件失败: [%s]`, n, err.Error())})
		}

		if _, err := ent_client.Template.Delete().Where(template.NameEQ(n)).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`删除模板"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除模板"%s"失败: [%s]`, n, err.Error())})
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`删除模板"%s"成功`, n)})
	}
	return res
}
func EditTemplate(template_msg model.Template, ent_client *ent.Client, logger *zap.Logger) error {
	template_instance := ent_client.Template.Update()
	if err := template_msg.EditProviders(); err != nil {
		return err
	}
	template_msg.EditRulesets()
	template_data, err := ent_client.Template.Query().Where(template.NameEQ(template_msg.Name)).First(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查找模板"%s"失败: [%s]`, template_msg.Name, err.Error()))
		return fmt.Errorf(`查找模板"%s"失败: [%s]`, template_msg.Name, err.Error())
	}
	providers_origin_map := make(map[string]bool)
	providers_update_map := make(map[string]bool)
	ruleset_orgin_map := make(map[string]bool)
	ruleset_update_map := make(map[string]bool)
	for _, r := range template_data.Route.Rule_sets {
		ruleset_orgin_map[r.Tag] = true
	}
	for _, r := range template_msg.Route.Rule_sets {
		ruleset_update_map[r.Tag] = true
	}
	for _, n := range template_data.Providers {
		providers_origin_map[n] = true
	}
	for _, n := range template_msg.Providers {
		providers_update_map[n] = true
	}
	for k := range providers_origin_map {
		if _, exist := providers_update_map[k]; !exist {
			provider_data, err := ent_client.Provider.Query().Where(provider.NameEQ(k)).First(context.Background())
			if err != nil {
				logger.Error(fmt.Sprintf(`查找模板"%s"的机场"%s"失败: [%s]`, template_msg.Name, k, err.Error()))
				return fmt.Errorf(`查找模板"%s"的机场"%s"失败: [%s]`, template_msg.Name, k, err.Error())
			}
			template_list := []string{}
			for _, v := range provider_data.Templates {
				if v != template_msg.Name {
					template_list = append(template_list, v)
				}
			}
			if _, err := ent_client.Provider.Update().Where(provider.NameEQ(k)).SetTemplates(template_list).Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf(`机场"%s"关联模板"%s"失败: [%s]`, k, template_msg.Name, err.Error()))
				return fmt.Errorf(`机场"%s"关联模板"%s"失败: [%s]`, k, template_msg.Name, err.Error())
			}
		}
	}
	for k := range providers_update_map {
		if _, exist := providers_origin_map[k]; !exist {
			provider_data, err := ent_client.Provider.Query().Where(provider.NameEQ(k)).First(context.Background())
			if err != nil {
				logger.Error(fmt.Sprintf(`查找模板"%s"的机场"%s"失败: [%s]`, template_msg.Name, k, err.Error()))
				return fmt.Errorf(`查找模板"%s"的机场"%s"失败: [%s]`, template_msg.Name, k, err.Error())
			}
			template_list := append(provider_data.Templates, template_msg.Name)
			if _, err := ent_client.Provider.Update().Where(provider.NameEQ(k)).SetTemplates(template_list).Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf(`机场"%s"关联模板"%s"失败: [%s]`, k, template_msg.Name, err.Error()))
				return fmt.Errorf(`机场"%s"关联模板"%s"失败: [%s]`, k, template_msg.Name, err.Error())
			}
		}
	}
	for k := range ruleset_orgin_map {
		if _, exist := ruleset_update_map[k]; !exist {
			ruleset_data, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(k)).First(context.Background())
			if err != nil {
				logger.Error(fmt.Sprintf(`查找模板"%s"的规则集"%s"失败: [%s]`, template_msg.Name, k, err.Error()))
				return fmt.Errorf(`查找模板"%s"的规则集"%s"失败: [%s]`, template_msg.Name, k, err.Error())
			}
			template_list := []string{}
			for _, v := range ruleset_data.Templates {
				if v != template_msg.Name {
					template_list = append(template_list, v)
				}
			}
			if _, err := ent_client.Ruleset.Update().Where(ruleset.NameEQ(k)).SetTemplates(template_list).Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf(`规则集"%s"关联模板"%s"失败: [%s]`, k, template_msg.Name, err.Error()))
				return fmt.Errorf(`规则集"%s"关联模板"%s"失败: [%s]`, k, template_msg.Name, err.Error())
			}
		}
	}
	for k := range ruleset_update_map {
		if _, exist := ruleset_orgin_map[k]; !exist {
			ruleset_data, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(k)).First(context.Background())
			if err != nil {
				logger.Error(fmt.Sprintf(`查找模板"%s"的规则集"%s"失败: [%s]`, template_msg.Name, k, err.Error()))
				return fmt.Errorf(`查找模板"%s"的规则集"%s"失败: [%s]`, template_msg.Name, k, err.Error())
			}
			template_list := append(ruleset_data.Templates, template_msg.Name)
			if _, err := ent_client.Ruleset.Update().Where(ruleset.NameEQ(k)).SetTemplates(template_list).Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf(`规则集"%s"关联模板"%s"失败: [%s]`, k, template_msg.Name, err.Error()))
				return fmt.Errorf(`规则集"%s"关联模板"%s"失败: [%s]`, k, template_msg.Name, err.Error())
			}
		}
	}
	template_msg.UpdateFillFields(template_instance)
	if err := template_instance.Where(template.NameEQ(template_msg.Name)).SetUpdated(true).SetInbounds(template_msg.Inbounds).SetOutboundGroups(template_msg.OutboundsGroup).SetProviders(template_msg.Providers).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf(`修改模板"%s"失败: [%s]`, template_msg.Name, err.Error()))
		return fmt.Errorf(`修改模板"%s"失败: [%s]`, template_msg.Name, err.Error())
	}
	return nil
}
