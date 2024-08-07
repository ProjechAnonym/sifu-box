package singbox

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"

	"github.com/huandu/go-clone"
)

// merge 函数负责将服务配置模板与一组服务提供者合并,生成最终的配置文件
// 它接受模板名称、项目目录、模板对象、模式标志、服务提供者列表和服务映射作为参数,
// 并返回一个错误切片,包含合并过程中出现的所有错误
func merge(key, projectDir string, template models.Template, mode bool, providers []models.Provider, serviceMap map[string][]models.Ruleset) []error {
    // 初始化同步等待组,用于等待所有并发的合并操作完成
    var jobs sync.WaitGroup
    // 根据服务提供者的数量创建错误通道,用于收集合并过程中出现的错误
    errChannel := make(chan error, len(providers))

    // 遍历服务提供者列表,为每个提供者启动一个并发的合并操作
    for _, provider := range providers {
        jobs.Add(1)
        go func() {
            defer jobs.Done()

            // 克隆模板对象,以避免在并发操作中修改原始模板
            tempTemplate := clone.Clone(template).(models.Template)
            // 为当前模板设置规则集
            tempTemplate.Route.Rule_set = append(tempTemplate.Route.Rule_set, SetRulesets(serviceMap)...)

            // 为当前模板和提供者设置规则
            tempTemplate.Route.Rules = append(tempTemplate.Route.Rules, SetRules(serviceMap, provider)...)

            // 为当前模板设置DNS规则
            tempTemplate.Dns.Rules = append(tempTemplate.Dns.Rules, SetDnsRules(serviceMap)...)

            // 合并出站配置,如果出错则记录错误并通知
            proxies, err := MergeOutbound(provider, serviceMap, tempTemplate.CustomOutbounds)
            if err != nil {
                utils.LoggerCaller(fmt.Sprintf("模板'%s'与'%s'节点合并失败", key, provider.Name), err, 1)
                errChannel <- fmt.Errorf("模板'%s'与'%s'节点合并失败", key, provider.Name)
                return
            }

            // 将合并后的出站配置添加到模板中
            tempTemplate.Outbounds = append(tempTemplate.Outbounds, proxies...)

            // 将模板对象序列化为JSON格式
            json, err := json.MarshalIndent(tempTemplate, "", "  ")
            if err != nil {
                utils.LoggerCaller("json序列化失败", err, 1)
                errChannel <- fmt.Errorf("模板'%s'与'%s'节点合并失败", key, provider.Name)
                return
            }

            // 根据模式标志决定标签名称
            var label string
            if mode {
                // 使用MD5加密提供者名称作为标签
                md5label, err := utils.EncryptionMd5(provider.Name)
                if err != nil {
                    utils.LoggerCaller("md5加密失败", err, 1)
                    errChannel <- fmt.Errorf("模板'%s'与'%s'节点合并失败", key, provider.Name)
                    return
                }
                label = md5label
            } else {
                label = provider.Name
            }

            // 将生成的JSON写入文件
            if err := utils.FileWrite(json, filepath.Join(projectDir, "static", key, fmt.Sprintf("%s.json", label))); err != nil {
                utils.LoggerCaller("写入文件失败", err, 1)
                errChannel <- fmt.Errorf("模板'%s'与'%s'节点合并失败", key, provider.Name)
                return
            }

            // 记录合并成功日志
            utils.LoggerCaller(fmt.Sprintf("模板'%s'与'%s'节点合并成功", key, provider.Name), nil, 1)
        }()
    }

    // 等待所有并发的合并操作完成
    jobs.Wait()
    // 关闭错误通道
    close(errChannel)

    // 收集合并过程中出现的所有错误
    var errs []error
    for err := range errChannel {
        errs = append(errs, err)
    }
    // 返回错误切片
    return errs
}
// Workflow 执行工作流,根据提供的特定参数处理模板和提供者
// 参数 specific 用于指定特定的提供者,如果没有提供,则处理所有提供者
// 返回一个错误切片,包含处理过程中遇到的所有错误
func Workflow(specific ...int) []error {
    // 获取项目目录
	projectDir,err := utils.GetValue("project-dir")
	if err != nil {
		utils.LoggerCaller("获取项目目录失败",err,1)
		return []error{fmt.Errorf("获取项目目录失败")}
	}
	// 获取模板配置
	templates, err := utils.GetValue("templates")
	if err != nil {
		utils.LoggerCaller("获取模板配置失败",err,1)
		return []error{fmt.Errorf("获取模板配置失败")}
	}
	// 初始化提供者切片
	var providers []models.Provider
	// 根据是否提供了特定参数,决定查询所有提供者还是指定的提供者
	if len(specific) == 0 {
		if err := utils.MemoryDb.Find(&providers).Error; err != nil {
			utils.LoggerCaller("获取provider配置失败",err,1)
			return []error{fmt.Errorf("获取provider配置失败")}
		}
	}else{
		if err := utils.MemoryDb.Find(&providers,specific).Error; err != nil {
			utils.LoggerCaller("获取provider配置失败",err,1)
			return []error{fmt.Errorf("获取provider配置失败")}
		}
	}
	// 检查提供者配置是否为空
	if len(providers) == 0 {
		utils.LoggerCaller("provider配置为空",nil,1)
		return []error{fmt.Errorf("provider配置为空")}
	}
	// 获取规则集配置
	var rulesets []models.Ruleset
	if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
		utils.LoggerCaller("获取ruleset配置失败",err,1)
		return []error{fmt.Errorf("获取ruleset配置失败")}
	}
	// 为提供者添加Clash标签并排序规则集
	newProviders,errs := AddClashTag(providers)
	newRulesets := SortRulesets(rulesets)
	// 初始化同步等待组和错误通道
	var workflow sync.WaitGroup
	errChannel := make(chan error,len(newProviders) * len(templates.(map[string]models.Template)))
	// 获取服务模式配置
	server,err := utils.GetValue("mode")
    if err != nil{
        utils.LoggerCaller("获取服务模式配置失败",err,1)
        return []error{fmt.Errorf("获取服务模式配置失败")}
    }
	mode := server.(models.Server).Mode
	// 遍历模板,为每个模板启动一个goroutine执行合并操作
	for key,value := range(templates.(map[string]models.Template)){
		workflow.Add(1)
		go func(){
			defer workflow.Done()
			errors := merge(key,projectDir.(string),value,mode,newProviders,newRulesets)
			for _,err := range errors{
				errChannel <- err
			}
		}()
	}
	// 等待所有goroutine完成
	workflow.Wait()
	// 关闭错误通道
	close(errChannel)
	// 收集所有错误
	for err := range errChannel{
		errs = append(errs,err)
	}
	// 返回错误切片
	return errs
}