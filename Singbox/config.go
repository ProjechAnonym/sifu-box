package singbox

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"path/filepath"
	utils "sifu-box/Utils"
	"sync"

	"github.com/bitly/go-simplejson"
	"github.com/huandu/go-clone"
)

// format_url 根据配置文件中的Proxy部分获取URL列表，并确保每个URL都包含"flag=clash"参数。
// 如果URL已经包含此参数，则不做更改；否则，添加该参数。
// 返回处理后的URL列表以及可能出现的错误。
func format_url() ([]string,error) {
    // 从配置文件中获取URL列表
    urls,err := utils.Get_value("Proxy","url")
    if err != nil {
        // 记录获取URL列表时的错误
        utils.Logger_caller("Get urls failed!",err,1)
        return nil,err
    }
    // 检查URL列表是否为空
    if len(urls.([]interface{})) == 0 {
        // 如果URL列表为空，创建并返回一个新的空列表和一个错误
        err = errors.New("get url list failed")
        utils.Logger_caller("Get urls failed!",err,1)
        return nil,err
    }
    // 初始化一个用于存储处理后URL的切片
    links := make([]string,len(urls.([]interface{})))
    // 初始化一个标志，用于指示是否找到了带有"flag=clash"的URL
    clash_tag := false
    // 遍历URL列表，处理每个URL
    for i, link := range urls.([]interface{}) {
        // 解析URL字符串
        parsed_url,err := url.Parse(link.(string))
        if err != nil {
            // 记录URL解析失败的错误
            utils.Logger_caller("Parse url failed!",err,1)
            return nil,err
        }
        // 获取URL的查询参数
        params := parsed_url.Query()
        // 遍历查询参数，检查是否已有"flag=clash"
        for key, values := range params {
            if key == "flag" && values[0] == "clash"{
                // 如果已存在"flag=clash"，将该URL添加到结果列表中，并设置标志
                links[i] = parsed_url.String()
                clash_tag = true
                break
            }
        }
        // 如果没有找到"flag=clash"，添加该参数并更新URL
        if !clash_tag{
            params.Add("flag","clash")
            parsed_url.RawQuery = params.Encode()
            links[i] = parsed_url.String()
        }
    }
    // 返回处理后的URL列表和nil错误
    return links,nil
}
func config_merge(template string) error{
	
	
	// 获取固定信息
	log,err := utils.Get_value(template,"log")
	if err != nil{
		utils.Logger_caller("Get log failed!",err,1)
		return err
	}
	dns,err := utils.Get_value(template,"dns")
	if err != nil{
		utils.Logger_caller("Get dns failed!",err,1)
		return err
	}
	inbounds,err := utils.Get_value(template,"inbounds")
	if err != nil{
		utils.Logger_caller("Get inbounds failed!",err,1)
		return err
	}
	experimental,err := utils.Get_value(template,"experimental")
	if err != nil{
		utils.Logger_caller("Get experimental failed!",err,1)
		return err
	}
	route,err := Merge_route(template)
	if err != nil{
		utils.Logger_caller("Get route failed!",err,1)
		return err
	}
	config := simplejson.New()
	config.Set("log", log)
	config.Set("dns", dns)
	config.Set("inbounds", inbounds)
	config.Set("route", route)
	config.Set("experimental", experimental)
	links,err := format_url()
	error_channel := make(chan error,len(links))
	var jobs sync.WaitGroup
	jobs.Add(len(links))
	if err != nil{
		return err
	}
	for i,link := range links{
		go func(link string,template string,config *simplejson.Json,index int) {
			defer jobs.Done()
			project_dir,err := utils.Get_value("project-dir")
			if err != nil{
				utils.Logger_caller("Get project dir failed!",err,1)
				error_channel <- fmt.Errorf("generate the %dth url config failed",index)
				return
			}
			full_config := clone.Clone(config)
			full_config.(*simplejson.Json).Set("Proxy", link)
			config_bytes,_ := full_config.(*simplejson.Json).EncodePretty()
			if err = utils.File_write(config_bytes,filepath.Join(project_dir.(string),"static",template,fmt.Sprintf("%d.json",index+1)),[]fs.FileMode{0666,0777});err != nil{
				utils.Logger_caller("Write config file failed!",err,1)
				error_channel <- fmt.Errorf("generate the %dth url config failed",index)
				return
			}
		}(link,template,config,i)
	}
	jobs.Wait()
	close(error_channel)
	for err := range error_channel{
		fmt.Println(err)
	}
	return nil
}
func Config_workflow(template string) error {
	if err := utils.Load_template(template); err != nil {
		utils.Logger_caller("load the template failed",err,1)
		return fmt.Errorf("load the %s template failed",template)
	}
	if err := utils.Load_config("Proxy"); err != nil {
		utils.Logger_caller("load the Proxy config failed",err,1)
		return fmt.Errorf("load the Proxy config failed")
	}
	config_merge(template)
	return nil
}
