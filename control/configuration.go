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
	"sifu-box/singbox"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FetchItems 用于从数据库中获取 Provider、Ruleset 和 Template 的列表, 并将结果封装为 []gin.H 返回
// 参数:
//   - ent_client: ent.Client 的实例, 用于数据库查询操作
//   - logger: zap.Logger 的实例, 用于记录日志信息
//
// 返回值:
//   - []gin.H: 包含三类数据（provider、ruleset、template）的响应列表, 每个元素是一个 gin.H,
//     其中包含 status、message 和 type 字段
func FetchItems(ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}

	// 查询所有 Provider 数据
	providers, err := ent_client.Provider.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取机场列表失败: [%s]", err.Error()))
		res = append(res, gin.H{"status": false, "message": fmt.Sprintf("获取机场列表失败: [%s]", err.Error())})
		return res
	}

	// 构造返回结构体列表
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

	// 填充 Provider 列表
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

	// 查询所有 Ruleset 数据
	rulesets, err := ent_client.Ruleset.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取规则集列表失败: [%s]", err.Error()))
		res = append(res, gin.H{"status": false, "message": fmt.Sprintf("获取规则集列表失败: [%s]", err.Error())})
		return res
	}

	// 填充 Ruleset 列表
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

	// 查询所有 Template 数据
	templates, err := ent_client.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取模板列表失败: [%s]", err.Error()))
		res = append(res, gin.H{"status": false, "message": fmt.Sprintf("获取模板列表失败: [%s]", err.Error())})
		return res
	}

	// 填充 Template 列表
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

// AddProvider 批量添加机场信息到数据库
// 参数:
//
//	providers []model.Provider - 需要添加的机场列表
//	ent_client *ent.Client - 数据库客户端实例
//	logger *zap.Logger - 日志记录器实例
//
// 返回值:
//
//	[]gin.H - 每个机场添加结果的响应信息列表, 包含状态和消息
func AddProvider(providers []model.Provider, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	// 遍历所有机场, 逐个添加到数据库
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

// AddRuleset 批量添加规则集到数据库
// 参数:
//   - rulesets: 规则集模型切片, 包含要添加的规则集信息
//   - ent_client: ent数据库客户端, 用于执行数据库操作
//   - logger: zap日志记录器, 用于记录操作日志
//
// 返回值:
//   - []gin.H: 操作结果切片, 每个元素包含状态和消息信息
func AddRuleset(rulesets []model.Ruleset, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	// 遍历规则集列表, 逐个添加到数据库
	for _, ruleset := range rulesets {
		// 构造并执行数据库插入操作, 如果失败则记录错误日志
		if err := ent_client.Ruleset.Create().SetName(ruleset.Name).SetPath(ruleset.Path).SetRemote(ruleset.Remote).SetBinary(ruleset.Binary).SetDownloadDetour(ruleset.DownloadDetour).SetUpdateInterval(ruleset.UpdateInterval).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加规则集"%s"失败: [%s]`, ruleset.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加规则集"%s"失败: [%s]`, ruleset.Name, err.Error())})
		} else {
			res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加规则集"%s"成功`, ruleset.Name)})
		}
	}
	return res
}

// EditProvider 编辑指定名称的机场信息
// name: 机场名称, 用于定位要编辑的机场记录
// path: 机场配置文件路径
// remote: 是否为远程机场标识
// ent_client: 数据库客户端实例, 用于执行数据库操作
// logger: 日志记录器实例, 用于记录操作日志
// 返回值: 操作成功返回nil, 失败返回相应的错误信息
func EditProvider(name, path string, remote bool, ent_client *ent.Client, logger *zap.Logger) error {
	// 检查指定名称的机场是否存在
	exist, err := ent_client.Provider.Query().Where(provider.NameEQ(name)).Exist(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查找机场"%s"失败: [%s]`, name, err.Error()))
		return fmt.Errorf(`查找机场"%s"失败: [%s]`, name, err.Error())
	} else if !exist {
		logger.Error(fmt.Sprintf(`未找到机场"%s"`, name))
		return fmt.Errorf(`机场"%s"不存在`, name)
	}

	// 更新机场信息
	if err := ent_client.Provider.Update().Where(provider.NameEQ(name)).SetPath(path).SetRemote(remote).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf(`修改机场"%s"失败: [%s]`, name, err.Error()))
	}
	return nil
}

// EditRuleset 编辑规则集信息
// name: 规则集名称, 用于标识和查询规则集
// path: 规则集路径, 指定规则集文件的存储路径
// update_interval: 更新间隔, 设置规则集的自动更新频率
// download_detour: 下载绕行设置, 指定下载时使用的网络策略
// remote: 是否为远程规则集, true表示从远程获取, false表示本地规则集
// binary: 是否为二进制格式, true表示二进制规则集, false表示文本格式
// ent_client: 数据库客户端, 用于执行数据库查询和更新操作
// logger: 日志记录器, 用于记录操作日志和错误信息
// 返回值: 操作成功返回nil, 失败返回相应的错误信息
func EditRuleset(name, path, update_interval, download_detour string, remote, binary bool, ent_client *ent.Client, logger *zap.Logger) error {
	// 检查规则集是否存在
	exist, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(name)).Exist(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查找规则集"%s"失败: [%s]`, name, err.Error()))
		return fmt.Errorf(`查找规则集"%s"失败: [%s]`, name, err.Error())
	} else if !exist {
		logger.Error(fmt.Sprintf(`未找到规则集"%s"`, name))
		return fmt.Errorf(`规则集"%s"不存在`, name)
	}

	// 更新规则集信息
	if err := ent_client.Ruleset.Update().Where(ruleset.NameEQ(name)).SetPath(path).SetDownloadDetour(download_detour).SetUpdateInterval(update_interval).SetBinary(binary).SetRemote(remote).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf(`修改机场"%s"失败: [%s]`, name, err.Error()))
	}
	return nil
}

// DeleteProvider 删除指定名称的机场提供商, 并更新相关模板中的引用
// 参数:
//   - name: 要删除的机场提供商名称列表
//   - ent_client: 数据库客户端实例, 用于查询和操作数据
//   - logger: 日志记录器, 用于输出错误日志
//
// 返回值:
//   - []gin.H: 每个元素表示一个操作结果, 包含状态和消息
func DeleteProvider(name []string, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	delete_providers_map := map[string]bool{}
	for _, n := range name {
		delete_providers_map[n] = true
	}

	// 遍历所有要删除的机场名称
	for _, n := range name {
		// 查询要删除的机场提供商信息
		provider_msg, err := ent_client.Provider.Query().Where(provider.NameEQ(n)).Select(provider.FieldName, provider.FieldPath, provider.FieldRemote, provider.FieldTemplates).First(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf(`查找机场"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找机场"%s"失败: [%s]`, n, err.Error())})
			continue
		}

		exit_status := false

		// 处理与该机场关联的所有模板, 移除对该机场的引用
		for _, template_name := range provider_msg.Templates {
			template_instance, err := ent_client.Template.Query().Where(template.NameEQ(template_name)).Select(template.FieldProviders, template.FieldOutboundGroups).First(context.Background())
			if err != nil {
				logger.Error(fmt.Sprintf(`查找模板"%s"失败: [%s]`, template_name, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找模板"%s"失败: [%s]`, template_name, err.Error())})
				exit_status = true
				break
			}
			template_msg := model.Template{OutboundsGroup: template_instance.OutboundGroups, Providers: template_instance.Providers}
			provider_list := []string{}
			for _, provider_name := range template_msg.Providers {
				if !delete_providers_map[provider_name] {
					provider_list = append(provider_list, provider_name)
				}
			}
			template_msg.Providers = provider_list
			if err := template_msg.EditProviders(); err != nil {
				logger.Error(fmt.Sprintf(`修改模板"%s"失败: [%s]`, template_name, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`修改模板"%s"失败: [%s]`, template_name, err.Error())})
				exit_status = true
				break
			}
			if err := ent_client.Template.Update().Where(template.NameEQ(template_name)).SetProviders(template_msg.Providers).SetOutboundGroups(template_msg.OutboundsGroup).SetUpdated(true).Exec(context.Background()); err != nil {
				logger.Error(fmt.Sprintf(`修改模板"%s"机场列表失败: [%s]`, template_name, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`修改模板"%s"机场列表失败: [%s]`, template_name, err.Error())})
				exit_status = true
				break
			}
		}
		if exit_status {
			continue
		}

		// 如果不是远程机场, 则删除本地文件
		if !provider_msg.Remote {
			if err := os.Remove(provider_msg.Path); err != nil {
				logger.Error(fmt.Sprintf(`删除机场"%s"文件失败: [%s]`, n, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除机场"%s"文件失败: [%s]`, n, err.Error())})
				continue
			}
		}

		// 从数据库中删除机场记录
		if _, err := ent_client.Provider.Delete().Where(provider.NameEQ(n)).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`删除机场"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除机场"%s"失败: [%s]`, n, err.Error())})
		}

		// 记录删除成功的结果
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`删除机场"%s"成功`, n)})
	}

	return res
}

// DeleteRuleset 用于删除指定名称的规则集, 并更新引用该规则集的模板配置
// 同时会根据规则集是否为本地文件决定是否删除对应文件, 并从数据库中移除规则集记录
//
// 参数:
//   - name: 要删除的规则集名称列表
//   - ent_client: 数据库客户端实例, 用于查询和操作规则集、模板等数据
//   - logger: 日志记录器, 用于记录错误日志
//
// 返回值:
//   - []gin.H: 每个元素表示一个规则集的删除结果, 包含状态和消息
func DeleteRuleset(name []string, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	delete_rulesets_map := map[string]bool{}
	for _, n := range name {
		delete_rulesets_map[n] = true
	}

	// 遍历所有要删除的规则集名称
	for _, n := range name {
		// 查询要删除的规则集信息
		rule_set, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(n)).Select(ruleset.FieldPath, ruleset.FieldRemote, ruleset.FieldTemplates).First(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf(`查找规则集"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找规则集"%s"失败: [%s]`, n, err.Error())})
			continue
		}

		exit_status := false

		// 更新引用当前规则集的所有模板中的规则集列表
		for _, template_name := range rule_set.Templates {
			template_instance, err := ent_client.Template.Query().Where(template.NameEQ(template_name)).Select(template.FieldRoute, template.FieldDNS).First(context.Background())
			if err != nil {
				logger.Error(fmt.Sprintf(`查找模板"%s"失败: [%s]`, template_name, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找模板"%s"失败: [%s]`, template_name, err.Error())})
				exit_status = true
				break
			}
			template_msg := model.Template{Route: &template_instance.Route, DNS: &template_instance.DNS}
			rule_set_list := []singbox.Rule_set{}
			for _, rule_set_msg := range template_msg.Route.Rule_sets {
				if !delete_rulesets_map[rule_set_msg.Tag] {
					rule_set_list = append(rule_set_list, rule_set_msg)
				}
			}
			template_msg.Route.Rule_sets = rule_set_list
			template_msg.EditRulesets()
			if err := ent_client.Template.Update().Where(template.NameEQ(template_name)).SetRoute(*template_msg.Route).SetDNS(*template_msg.DNS).SetUpdated(true).Exec(context.Background()); err != nil {
				logger.Error(fmt.Sprintf(`修改模板"%s"机场列表失败: [%s]`, template_name, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`修改模板"%s"规则集列表失败: [%s]`, template_name, err.Error())})
				exit_status = true
				break
			}
		}
		if exit_status {
			continue
		}
		// 如果是本地规则集, 删除对应的文件
		if !rule_set.Remote {
			if err := os.Remove(rule_set.Path); err != nil {
				logger.Error(fmt.Sprintf(`删除规则集"%s"文件失败: [%s]`, n, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除规则集"%s"文件失败: [%s]`, n, err.Error())})
				continue
			}
		}
		// 从数据库中删除规则集记录
		if _, err := ent_client.Ruleset.Delete().Where(ruleset.NameEQ(n)).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`删除规则集"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除规则集"%s"失败: [%s]`, n, err.Error())})
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`删除规则集"%s"成功`, n)})
	}
	return res
}

// AddTemplate 向数据库中添加一个新的模板记录
// template: 要添加的模板对象, 包含模板的所有信息
// ent_client: 数据库客户端, 用于执行数据库操作
// logger: 日志记录器, 用于记录错误日志
// 返回值: 如果添加成功返回nil, 否则返回相应的错误信息
func AddTemplate(template model.Template, ent_client *ent.Client, logger *zap.Logger) error {
	// 创建数据库模板记录构建器
	ent_template := ent_client.Template.Create()
	template.CreateFillFields(ent_template)

	// 关联提供商表数据
	if err := template.LinkProvidersTable(ent_client); err != nil {
		logger.Error(err.Error())
		return err
	}

	// 关联规则集表数据
	if err := template.LinkRulesetsTable(ent_client); err != nil {
		logger.Error(err.Error())
		return err
	}

	// 编辑机场数据
	if err := template.EditProviders(); err != nil {
		return err
	}

	// 编辑规则集数据
	template.EditRulesets()

	// 设置模板字段并执行数据库插入操作
	if err := ent_template.SetUpdated(true).SetName(template.Name).SetProviders(template.Providers).SetInbounds(template.Inbounds).SetOutboundGroups(template.OutboundsGroup).Exec(context.Background()); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// DeleteTemplate 用于删除指定名称的模板, 包括从数据库中删除模板记录、解除关联关系以及删除对应的配置文件
// 参数:
//   - name: 要删除的模板名称列表
//   - work_dir: 工作目录路径, 用于定位模板配置文件
//   - ent_client: ent 数据库客户端, 用于执行数据库操作
//   - logger: zap 日志记录器, 用于记录错误和操作日志
//
// 返回值:
//   - []gin.H: 每个模板的删除结果, 包含状态和消息
func DeleteTemplate(name []string, work_dir string, ent_client *ent.Client, logger *zap.Logger) []gin.H {
	res := []gin.H{}

	// 遍历所有要删除的模板名称
	for _, n := range name {
		// 查询模板信息（名称、提供商、路由）
		template_msg, err := ent_client.Template.Query().Where(template.NameEQ(n)).Select(template.FieldName, template.FieldProviders, template.FieldRoute).First(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf(`查找模板"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找模板"%s"失败: [%s]`, n, err.Error())})
			continue
		}

		// 构造模板实例, 用于解除与 Providers 和 Rulesets 表的关联
		template_instance := model.Template{Route: &template_msg.Route, Providers: template_msg.Providers, Name: template_msg.Name}

		// 解除与 Providers 表的关联
		if err := template_instance.UnLinkProvidersTable(ent_client); err != nil {
			logger.Error(err.Error())
			res = append(res, gin.H{"status": false, "message": err.Error()})
			continue
		}

		// 解除与 Rulesets 表的关联
		if err := template_instance.UnLinkRulesetsTable(ent_client); err != nil {
			logger.Error(err.Error())
			res = append(res, gin.H{"status": false, "message": err.Error()})
			continue
		}

		// 删除模板对应的配置文件（如果存在）
		configFilePath := filepath.Join(work_dir, "sing-box", "config", fmt.Sprintf(`%s.json`, fmt.Sprintf(`%x`, md5.Sum([]byte(template_msg.Name)))))
		if _, err := os.Stat(configFilePath); err == nil {
			if err := os.Remove(configFilePath); err != nil {
				logger.Error(fmt.Sprintf(`删除模板"%s"配置文件失败: [%s]`, n, err.Error()))
				res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除模板"%s"配置文件失败: [%s]`, n, err.Error())})
				continue
			}
		} else if !os.IsNotExist(err) {
			logger.Error(fmt.Sprintf(`查找模板"%s"配置文件失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`查找模板"%s"配置文件失败: [%s]`, n, err.Error())})
			continue
		}

		// 从数据库中删除模板记录
		if _, err := ent_client.Template.Delete().Where(template.NameEQ(n)).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`删除模板"%s"失败: [%s]`, n, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`删除模板"%s"失败: [%s]`, n, err.Error())})
			continue
		}

		// 添加成功删除的响应
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`删除模板"%s"成功`, n)})
	}

	return res
}

// EditTemplate 用于编辑指定模板的信息, 并同步更新其关联的 Providers 和 Rulesets
// 该函数会处理模板字段的更新、关联关系的增删, 并记录操作日志
//
// 参数:
//   - template_msg: 包含待更新模板信息的 model.Template 实例
//   - ent_client: Ent ORM 客户端实例, 用于数据库操作
//   - logger: Zap 日志记录器实例, 用于记录错误日志
//
// 返回值:
//   - error: 如果在执行过程中发生错误, 则返回相应的错误信息；否则返回 nil
func EditTemplate(template_msg model.Template, ent_client *ent.Client, logger *zap.Logger) error {
	// 初始化模板更新实例
	template_instance := ent_client.Template.Update()

	// 调用模板对象的方法更新 Providers 和 Rulesets 的相关数据
	if err := template_msg.EditProviders(); err != nil {
		return err
	}
	template_msg.EditRulesets()

	// 查询当前数据库中同名模板的数据
	template_data, err := ent_client.Template.Query().Where(template.NameEQ(template_msg.Name)).First(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查找模板"%s"失败: [%s]`, template_msg.Name, err.Error()))
		return fmt.Errorf(`查找模板"%s"失败: [%s]`, template_msg.Name, err.Error())
	}

	// 构建原始和更新后的 Providers 与 Rulesets 映射, 用于后续比较差异
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

	// 处理被移除的 Providers：从对应 Provider 的 Templates 列表中删除当前模板名称
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

	// 处理新增的 Providers：将当前模板名称添加到对应 Provider 的 Templates 列表中
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

	// 处理被移除的 Rulesets：从对应 Ruleset 的 Templates 列表中删除当前模板名称
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

	// 处理新增的 Rulesets：将当前模板名称添加到对应 Ruleset 的 Templates 列表中
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

	// 填充模板更新字段并执行数据库更新操作
	template_msg.UpdateFillFields(template_instance)
	if err := template_instance.Where(template.NameEQ(template_msg.Name)).SetUpdated(true).SetInbounds(template_msg.Inbounds).SetOutboundGroups(template_msg.OutboundsGroup).SetProviders(template_msg.Providers).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf(`修改模板"%s"失败: [%s]`, template_msg.Name, err.Error()))
		return fmt.Errorf(`修改模板"%s"失败: [%s]`, template_msg.Name, err.Error())
	}

	return nil
}
