package singbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sifu-box/models"
	"sync"

	"go.uber.org/zap"
)

func merge(providerList []models.Provider, rulesetsList []models.RuleSet, templates map[string]models.Template, logger *zap.Logger) []error{
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
			var err error
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
			for _, template := range templates {
				template.Dns.SetDNSRules(rulesetsList)
				template.Route.SetRuleSet(rulesetsList, logger)
				template.Route.SetRules(provider, rulesetsList, logger)
				template.SetOutbounds(outbounds) 
				a, _ := json.Marshal(template)
				os.WriteFile(fmt.Sprintf("E:\\MyProject\\sifu-box@1.1.0\\static\\%s.json",provider.Name),a,0666)
			}
			
		}()
	}
	jobs.Wait()	
	return nil
}

