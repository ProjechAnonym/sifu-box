package generate

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/nodes"

	"go.uber.org/zap"
)

func preProcess(ent_client *ent.Client, logger *zap.Logger) (map[string]bool, map[string][]map[string]any, []error, error) {
	original_providers, err := ent_client.Provider.Query().All(context.Background())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("获取机场信息失败: [%s]", err.Error())
	}
	fetch_errors := nodes.Fetch(original_providers, ent_client, logger)
	update_providers, err := ent_client.Provider.Query().Select(provider.FieldNodes, provider.FieldUUID, provider.FieldName).All(context.Background())
	if err != nil {
		return nil, nil, fetch_errors, fmt.Errorf("获取机场信息失败: [%s]", err.Error())
	}
	update_map := map[string]bool{}
	provider_nodes := map[string][]map[string]any{}
	for _, update_provider := range update_providers {
		data, err := json.Marshal(update_provider.Nodes)
		if err != nil {
			return nil, nil, fetch_errors, fmt.Errorf(`"%s"出错, 序列化出站节点失败: [%s]`, update_provider.Name, err.Error())
		}
		if update_provider.UUID != fmt.Sprintf("%x", md5.Sum([]byte(data))) {
			update_map[update_provider.Name] = true
		} else {
			update_map[update_provider.Name] = false
		}
		provider_nodes[update_provider.Name] = update_provider.Nodes
	}
	return update_map, provider_nodes, fetch_errors, nil
}
func Process(ent_client *ent.Client, logger *zap.Logger) []error {
	update_map, provider_nodes, fetch_errors, err := preProcess(ent_client, logger)
	if err != nil {
		logger.Error(fmt.Sprintf(`预处理出错: [%s]`, err.Error()))
	}
	templates, err := ent_client.Template.Query().All(context.Background())
	if err != nil {
		fetch_errors = append(fetch_errors, fmt.Errorf(`查询模板出错: [%s]`, err.Error()))
		return fetch_errors
	}
	for _, template := range templates {
		for _, name := range template.Providers {
			if update_map[name] {
				config := Config{}
				config.Generate(template, provider_nodes, logger)
			}
		}
	}
	return nil
}
