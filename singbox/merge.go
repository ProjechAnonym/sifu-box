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

func merge(providerList []models.Provider, rulesetsList []models.RuleSet, templates map[string]models.Template, workDir string, server bool, logger *zap.Logger) []error{
	providers, errors := formatProviderURL(providerList, logger)
	if errors != nil {
		return errors
	}
	requestClient := http.DefaultClient
	var jobs sync.WaitGroup
	var errChan = make(chan error, 5)
	var countChan = make(chan int, 5)
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
	jobs.Wait()	
	return errors
}

