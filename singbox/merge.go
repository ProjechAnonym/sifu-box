package singbox

import (
	"net/http"
	"sifu-box/models"
	"sync"

	"go.uber.org/zap"
)

func merge(providerList []models.Provider, logger *zap.Logger) []error{
	providers, errors := formatProviderURL(providerList, logger)
	if errors != nil {
		return errors
	}
	requestClient := http.DefaultClient
	var jobs sync.WaitGroup
	for _, provider := range providers {
		jobs.Add(1)
		go func(){
			defer jobs.Done()
			fetchProviderInfo(provider, requestClient, logger)
		}()
	}
	jobs.Wait()	
	return nil
}