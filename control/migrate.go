package control

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Export 用于从 ent 数据库和 buntdb 数据库中导出配置信息, 并将其序列化为 YAML 格式的字节数据
// 参数:
//   - ent_client: ent 数据库客户端, 用于查询 Provider、Ruleset 和 Template 数据
//   - bunt_client: buntdb 数据库客户端（当前未使用）
//   - logger: zap 日志记录器, 用于记录错误日志
//
// 返回值:
//   - []byte: 序列化后的 YAML 配置数据
//   - error: 如果在查询或序列化过程中发生错误, 则返回相应的错误信息
func Export(ent_client *ent.Client, bunt_client *buntdb.DB, logger *zap.Logger) ([]byte, error) {
	// 查询所有 Provider 数据
	providers, err := ent_client.Provider.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查询机场数据失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`查询机场数据失败: [%s]`, err.Error())
	}

	// 构造包含额外字段的 Provider 列表结构
	provider_list := []struct {
		model.Provider `json:",inline" yaml:",inline"`
		UUID           string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
		Nodes          []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
		Templates      []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{}

	// 构造 Ruleset 列表结构
	ruleset_list := []struct {
		model.Ruleset `json:",inline" yaml:",inline"`
		Templates     []string `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{}

	// 初始化 Template 列表
	template_list := []model.Template{}

	// 遍历并转换 Provider 数据
	for _, provider := range providers {
		provider_list = append(provider_list, struct {
			model.Provider `json:",inline" yaml:",inline"`
			UUID           string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
			Nodes          []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
			Templates      []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
		}{
			Provider: model.Provider{
				Name:   provider.Name,
				Path:   provider.Path,
				Remote: provider.Remote},
			Nodes:     provider.Nodes,
			UUID:      provider.UUID,
			Templates: provider.Templates,
		})
	}

	// 查询所有 Ruleset 数据
	rulesets, err := ent_client.Ruleset.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查询规则集数据失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`查询规则集数据失败: [%s]`, err.Error())
	}

	// 遍历并转换 Ruleset 数据
	for _, ruleset := range rulesets {
		ruleset_list = append(ruleset_list, struct {
			model.Ruleset `json:",inline" yaml:",inline"`
			Templates     []string `json:"templates,omitempty" yaml:"templates,omitempty"`
		}{
			Templates: ruleset.Templates,
			Ruleset: model.Ruleset{
				Name:           ruleset.Name,
				Path:           ruleset.Path,
				Remote:         ruleset.Remote,
				Binary:         ruleset.Binary,
				UpdateInterval: ruleset.UpdateInterval,
				DownloadDetour: ruleset.DownloadDetour,
			},
		})
	}

	// 查询所有 Template 数据
	templates, err := ent_client.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查询模板数据失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`查询模板数据失败: [%s]`, err.Error())
	}

	// 遍历并转换 Template 数据
	for _, template := range templates {
		template_list = append(template_list, model.Template{
			Name:           template.Name,
			Ntp:            &template.Ntp,
			Inbounds:       template.Inbounds,
			Providers:      template.Providers,
			Route:          &template.Route,
			OutboundsGroup: template.OutboundGroups,
			DNS:            &template.DNS,
			Experiment:     &template.Experiment,
			Log:            &template.Log,
		})
	}

	// 将所有数据结构合并并序列化为 YAML 格式
	content, err := yaml.Marshal(struct {
		Providers []struct {
			model.Provider `json:",inline" yaml:",inline"`
			UUID           string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
			Nodes          []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
			Templates      []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"providers,omitempty" yaml:"providers,omitempty"`
		Rulesets []struct {
			model.Ruleset `json:",inline" yaml:",inline"`
			Templates     []string `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"rulesets,omitempty" yaml:"rulesets,omitempty"`
		Templates []model.Template `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{Templates: template_list, Providers: provider_list, Rulesets: ruleset_list})
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf(`反序列化配置信息失败: [%s]`, err.Error())
	}

	// 返回序列化后的 YAML 数据
	return content, nil
}

// Import 用于将传入的配置内容解析并导入到数据库中, 支持 providers、rulesets 和 templates 的导入
// 参数:
//   - content: 配置文件的字节内容, 格式应为 YAML 或 JSON
//   - ent_client: 数据库客户端, 用于执行数据插入操作
//   - logger: 日志记录器, 用于记录错误和操作日志
//
// 返回值:
//   - []gin.H: 每个元素表示一个导入项的结果, 包含状态和消息
//   - error: 如果在解析配置时发生错误, 则返回错误信息
func Import(content []byte, ent_client *ent.Client, logger *zap.Logger) ([]gin.H, error) {
	// 定义结构体用于解析配置文件内容, 包含 Providers、Rulesets 和 Templates 三部分
	setting := struct {
		Providers []struct {
			model.Provider `json:",inline" yaml:",inline"`
			UUID           string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
			Nodes          []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
			Templates      []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"providers,omitempty" yaml:"providers,omitempty"`
		Rulesets []struct {
			model.Ruleset `json:",inline" yaml:",inline"`
			Templates     []string `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"rulesets,omitempty" yaml:"rulesets,omitempty"`
		Templates []model.Template `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{}

	// 将配置内容反序列化为结构体
	if err := yaml.Unmarshal(content, &setting); err != nil {
		logger.Error(fmt.Sprintf(`序列化配置文件失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`序列化配置文件失败: [%s]`, err.Error())
	}

	// 存储导入结果的列表
	res := []gin.H{}

	// 导入 Providers 数据
	for _, provider := range setting.Providers {
		if err := ent_client.Provider.Create().
			SetName(provider.Name).
			SetPath(provider.Path).
			SetRemote(provider.Remote).
			SetUUID(provider.UUID).
			SetNodes(provider.Nodes).
			SetTemplates(provider.Templates).
			Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加机场"%s"失败: [%s]`, provider.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加机场"%s"失败: [%s]`, provider.Name, err.Error())})
			continue
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加机场"%s"成功`, provider.Name)})
	}

	// 导入 Rulesets 数据
	for _, ruleset := range setting.Rulesets {
		if err := ent_client.Ruleset.Create().
			SetName(ruleset.Name).
			SetPath(ruleset.Path).
			SetRemote(ruleset.Remote).
			SetBinary(ruleset.Binary).
			SetUpdateInterval(ruleset.UpdateInterval).
			SetDownloadDetour(ruleset.DownloadDetour).
			SetTemplates(ruleset.Templates).
			Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加规则集%s失败: [%s]`, ruleset.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加规则集"%s"失败: [%s]`, ruleset.Name, err.Error())})
			continue
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加规则集"%s"成功`, ruleset.Name)})
	}

	// 导入 Templates 数据
	for _, template := range setting.Templates {
		template_instance := ent_client.Template.Create()
		template.CreateFillFields(template_instance)
		if err := template_instance.
			SetInbounds(template.Inbounds).
			SetName(template.Name).
			SetOutboundGroups(template.OutboundsGroup).
			SetProviders(template.Providers).
			SetUpdated(true).
			Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加模板"%s"失败: [%s]`, template.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加模板"%s"失败: [%s]`, template.Name, err.Error())})
			continue
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加模板"%s"成功`, template.Name)})
	}

	// 返回导入结果和 nil 错误
	return res, nil
}
