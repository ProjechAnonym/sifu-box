package application

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/nodes"
	"sync"

	"go.uber.org/zap"
)

// preProcess 预处理函数, 用于获取和处理机场信息
// 参数:
//   - ent_client: ent数据库客户端, 用于查询和操作数据
//   - logger: zap日志记录器, 用于记录日志信息
//
// 返回值:
//   - map[string]bool: 机场更新状态映射, key为机场名称, value为是否已更新
//   - map[string][]map[string]any: 机场节点信息映射, key为机场名称, value为节点信息列表
//   - []error: 获取过程中产生的错误列表
//   - error: 函数执行过程中的错误信息
func preProcess(ent_client *ent.Client, logger *zap.Logger) (map[string]bool, map[string][]map[string]any, []error, error) {
	// 查询所有机场信息
	original_providers, err := ent_client.Provider.Query().All(context.Background())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("获取机场信息失败: [%s]", err.Error())
	}

	// 获取机场信息并记录可能产生的错误
	fetch_errors := nodes.Fetch(original_providers, ent_client, logger)

	// 重新查询机场信息, 只选择需要的字段
	update_providers, err := ent_client.Provider.Query().Select(provider.FieldNodes, provider.FieldUpdated, provider.FieldName).All(context.Background())
	if err != nil {
		return nil, nil, fetch_errors, fmt.Errorf("获取机场信息失败: [%s]", err.Error())
	}

	// 构建机场更新状态和节点信息的映射
	update_map := map[string]bool{}
	provider_nodes := map[string][]map[string]any{}
	for _, update_provider := range update_providers {
		update_map[update_provider.Name] = update_provider.Updated
		provider_nodes[update_provider.Name] = update_provider.Nodes
	}

	return update_map, provider_nodes, fetch_errors, nil
}

// Process 处理指定目录下的模板配置文件生成任务
// 该函数会根据数据库中的模板信息和机场节点信息, 生成对应的配置文件
// 参数:
//   - dir: 指定生成配置文件的根目录路径
//   - ent_client: 数据库客户端, 用于查询模板和节点信息
//   - logger: 日志记录器, 用于记录处理过程中的日志信息
//
// 返回值:
//   - []error: 处理过程中发生的错误列表
func Process(dir string, ent_client *ent.Client, logger *zap.Logger) []error {
	// 预处理阶段：获取变化机场信息、机场节点信息以及出现的错误列表
	update_map, provider_nodes, fetch_errors, err := preProcess(ent_client, logger)
	if err != nil {
		logger.Error(fmt.Sprintf(`预处理出错: [%s]`, err.Error()))
	}

	// 查询所有模板信息
	templates, err := ent_client.Template.Query().All(context.Background())
	if err != nil {
		fetch_errors = append(fetch_errors, fmt.Errorf(`查询模板出错: [%s]`, err.Error()))
		return fetch_errors
	}

	// 创建并发控制相关变量
	var jobs sync.WaitGroup
	var errChan = make(chan error, 5)
	var countChan = make(chan int, 5)

	// 启动监控协程, 用于收集处理结果和错误信息
	jobs.Add(1)
	go func() {
		defer func() {
			jobs.Done()
			var ok bool
			if _, ok = <-countChan; ok {
				close(countChan)
			}
			if _, ok = <-errChan; ok {
				close(errChan)
			}
		}()

		// 等待所有模板处理完成
		sum := 0
		for {
			if sum == len(templates) {
				return
			}
			select {
			case count, ok := <-countChan:
				if !ok {
					return
				}
				sum += count
			case err, ok := <-errChan:
				if !ok {
					return
				}
				fetch_errors = append(fetch_errors, err)
			}
		}
	}()

	// 为每个模板启动独立的处理协程
	for _, template := range templates {
		jobs.Add(1)
		go func() {
			defer func() {
				jobs.Done()
				countChan <- 1
			}()

			// 遍历模板关联的机场, 检查是否需要更新
			for _, name := range template.Providers {

				if update_map[name] || template.Updated {
					config := Config{}

					// 生成配置文件
					if err := config.Generate(dir, template, provider_nodes, logger); err != nil {
						logger.Error(err.Error())
						errChan <- err
					}
					if err := ent_client.Template.UpdateOne(template).SetUpdated(false).Exec(context.Background()); err != nil {
						errChan <- err
					}
					break
				}
			}
		}()
	}

	// 等待所有协程处理完成
	jobs.Wait()
	return fetch_errors
}
