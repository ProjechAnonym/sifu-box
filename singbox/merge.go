package singbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"

	"go.uber.org/zap"
)

// merge函数合并并处理提供商列表、规则集和模板, 生成相应的配置文件
// 参数:
// - providerList: 一个Provider类型的切片, 包含提供商的信息
// - rulesetsList: 一个RuleSet类型的切片, 包含规则集信息
// - templates: 一个映射, 键为字符串, 值为Template类型, 包含模板信息
// - workDir: 字符串类型, 表示工作目录
// - server: 布尔类型, 指示是否为服务器模式
// - logger: *zap.Logger类型, 用于日志记录
// 返回值:
// - []error: 一个错误切片, 包含处理过程中可能发生的错误
func merge(providerList []models.Provider, rulesetsList []models.RuleSet, templates map[string]models.Template, workDir string, server bool, logger *zap.Logger) []error{
    // 格式化机场URL, 在URL中添加flag参数, 并设置flag参数为clash
    providers, errors := formatProviderURL(providerList, logger)
	// 返回所出现的错误
    if errors != nil {
        return errors
    }
    // 初始化HTTP客户端
    requestClient := http.DefaultClient
    // 初始化错误、计数通道以及多线程计数变量
    var jobs sync.WaitGroup
    var errChan = make(chan error, 5)
    var countChan = make(chan int, 5)
    // 启动监控 goroutine, 负责收集错误和计数
    jobs.Add(1)
    go func(){
        defer func(){
            jobs.Done()
            var ok bool
            if _, ok = <- countChan; ok {close(countChan)}
            if _, ok = <- errChan; ok {close(errChan)}
        }()
        sum := 0
        for {
            if sum == len(providers) {return}
            select {
                case count, ok := <- countChan:
                    if !ok {return}
                    sum += count
                case err,ok := <- errChan:
                    if !ok {return}
                    errors = append(errors, err)
            }
        }	
    }()
    // 遍历每个机场, 启动处理 goroutine
    for _, provider := range providers {
        jobs.Add(1)
        go func(){
            defer func(){
                jobs.Done()
                countChan <- 1
            }()
            var outbounds []models.Outbound
            var providerName string
            var err error
            // 根据是否为服务器模式决定是否计算机场名称的哈希码
            if server {
                providerName, err = utils.EncryptionMd5(provider.Name)
                if err != nil {
                    logger.Error(fmt.Sprintf("'%s'生成哈希码失败: [%s]", provider.Name, err.Error()))
                    errChan <- fmt.Errorf("'%s'出错: '%s'生成哈希码失败", provider.Name, err.Error())
                    return
                }
            }else{
                providerName = provider.Name
            }
            // 根据机场是否为远程决定获取outbounds的方式
            if provider.Remote {
                outbounds, err = fetchFromRemote(provider, requestClient, logger)
            }else{
                outbounds, err = fetchFromLocal(provider, logger)
            }
            if err != nil {
                logger.Error(fmt.Sprintf("获取'%s'的outbounds失败: [%s]", provider.Name, err.Error()))
                errChan <- err
                return
            }
            // 处理并添加自动测试延迟出站配置和选择指定出站配置
            tags := make([]string, len(outbounds))
            for i, outbound := range outbounds {
                tags[i] = outbound.GetTag()
            }
            outbounds, tags, err = addURLTestOutbound(outbounds, tags, logger)
            if err != nil {
                logger.Error(fmt.Sprintf("'%s'生成auto出站失败: [%s]", provider.Name, err.Error()))
                errChan <- fmt.Errorf("'%s'出错: %s", provider.Name, err.Error())
                return
            }
            outbounds, err = addSelectorOutbound(provider.Name, outbounds, rulesetsList, tags, logger)
            if err != nil {
                logger.Error(fmt.Sprintf("'%s'生成默认selector出站失败: [%s]", provider.Name, err.Error()))
            }
            // 遍历每个模板, 生成并写入配置文件
            for key, template := range templates {
                template.Dns.SetDNSRules(rulesetsList)
                template.Route.SetRuleSet(rulesetsList, logger)
                template.Route.SetRules(provider, rulesetsList, logger)
                template.SetOutbounds(outbounds) 
                singboxConfigByte, err := json.Marshal(template)
                if err != nil {
                    logger.Error(fmt.Sprintf("反序列化'%s'基于模板'%s'的配置文件失败: [%s]", provider.Name, key, err.Error()))
                    errChan <- fmt.Errorf("'%s'出错: 反序列化基于'%s'模板的配置文件失败", provider.Name, key)
                }
                
                if err := utils.WriteFile(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR, key, fmt.Sprintf("%s.json", providerName)), singboxConfigByte, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644); err != nil {
                    logger.Error(fmt.Sprintf("'%s'基于模板'%s'生成配置文件失败: [%s]", provider.Name, key, err.Error()))
                    errChan <- fmt.Errorf("'%s'出错: '%s'生成配置文件失败", provider.Name, key)
                }
            }
            
        }()
    }
    // 等待所有任务完成
    jobs.Wait()	
    return errors
}

