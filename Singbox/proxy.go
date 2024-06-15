package singbox

import (
	"errors"
	"fmt"
	"net/http"
	utils "sifu-box/Utils"

	"github.com/gocolly/colly/v2"
	"gopkg.in/yaml.v3"
)
func handle_yaml(config map[string]interface{},header *http.Header,host string,template string) ([]map[string]interface{},error){
	
	
	if len(config["proxies"].([]interface{})) == 0{
		utils.Logger_caller("no proxies",fmt.Errorf("config without proxies"),1)
		return nil,errors.New("no proxies")
	}
	proxies := make([]map[string]interface{},len(config["proxies"].([]interface{})))
	
	for i,proxy := range config["proxies"].([]interface{}){
		proxies[i],_ = Format_yaml(proxy.(map[string]interface{}),template)
	}
	return proxies,nil
}
func fetch_proxies(url string,template string)  ([]map[string]interface{},error){
	var proxies []map[string]interface{}
	var err error
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		content := map[string]interface{}{}
		if err := yaml.Unmarshal(r.Body,&content); err != nil{
			msg := fmt.Sprintf("Parse %s proxies yaml failed!",r.Request.URL.Host)
			utils.Logger_caller(msg,err,1)
		}else{

		}
		proxies,err = handle_yaml(content,r.Headers,r.Request.URL.Host,template)
	})
	c.OnError(func(r *colly.Response, e error) {
		utils.Logger_caller(fmt.Sprintf("Connect to %s failed!",r.Request.URL.Host),e,1)
		err = e
	})
	c.Visit(url)
	if err != nil{
		return nil,err
	}
	return proxies,nil
}

func Merge_outbounds(url string,template string) ([]map[string]interface{},error){
	proxies,err := fetch_proxies(url,template)
	if err != nil{
		utils.Logger_caller("fetch proies failed!",err,1)
		return nil,err
	}
	return proxies,err
}