package singbox

import (
	"fmt"
	"net/http"
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
			if err := addURLTestOutbound(outbounds, tags, logger); err != nil {
				logger.Error(fmt.Sprintf("'%s'生成auto出站失败: [%s]", provider.Name, err.Error()))
				errChan <- fmt.Errorf("'%s'出错: %s", provider.Name, err.Error())
				return
			}
			targets := filterRulesetList(rulesetsList)
			var selector models.Selector
			selectorMap := map[string]interface{}{"type": "selector", "interrupt_exist_connections": false, "outbounds": tags, "tag": "select"}
			outbound = &selector
			outbound, err = outbound.Transform(selectorMap, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("'%s'生成默认selector出站失败: [%s]", provider.Name, err.Error()))
			}
			outbounds = append(outbounds, outbound)
			for _, target := range targets {
				selectorMap["tag"] = target
				outbound = &selector
				outbound, err = outbound.Transform(selectorMap, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("'%s'生成%s出站失败: [%s]", provider.Name, target, err.Error()))
				}
				outbounds = append(outbounds, outbound)
			}
			
		}()
	}
	jobs.Wait()	
	fmt.Println(errors)
	return nil
}

