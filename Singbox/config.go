package singbox

import (
	"crypto/md5"
	"encoding/hex"
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
func encryption_md5(str string) (string,error) {
	h := md5.New()
	_,err := h.Write([]byte(str))
	if err != nil {
		return "",err
	}
	return hex.EncodeToString(h.Sum(nil)),nil
}
// format_url 根据配置文件中的Proxy部分获取URL列表,并确保每个URL都包含"flag=clash"参数
// 如果URL已经包含此参数,则不做更改；否则,添加该参数
// 返回处理后的URL列表以及可能出现的错误
type config_link struct {
	url string
	proxy bool
	label string
}
func format_url() ([]config_link,error) {

    // 从配置文件中获取URL列表
    urls,err := utils.Get_value("Proxy","url")
    if err != nil {
        // 记录获取URL列表时的错误
        utils.Logger_caller("Get urls failed!",err,1)
        return nil,err
    }
    // 检查URL列表是否为空
    if len(urls.([]interface{})) == 0 {
        // 如果URL列表为空,创建并返回一个新的空列表和一个错误
        err = errors.New("get url list failed")
        utils.Logger_caller("Get urls failed!",err,1)
        return nil,err
    }
    // 初始化一个用于存储处理后URL的切片
    links := make([]config_link,len(urls.([]interface{})))
    // 初始化一个标志,用于指示是否找到了带有"flag=clash"的URL
    clash_tag := false
    // 遍历URL列表,处理每个URL
    for i, link := range urls.([]interface{}) {
		// 从URL映射中提取并设置代理标志和标签
		links[i].proxy = link.(map[string]interface{})["proxy"].(bool)
		links[i].label = link.(map[string]interface{})["label"].(string)
        // 解析URL字符串
        parsed_url,err := url.Parse(link.(map[string]interface{})["url"].(string))
        if err != nil {
            // 记录URL解析失败的错误
            utils.Logger_caller("Parse url failed!",err,1)
            return nil,err
        }
        // 获取URL的查询参数
        params := parsed_url.Query()
        // 遍历查询参数,检查是否已有"flag=clash"
        for key, values := range params {
            if key == "flag" && values[0] == "clash"{
                // 如果已存在"flag=clash",将该URL添加到结果列表中,并设置标志
                links[i].url = parsed_url.String()
                clash_tag = true
                break
            }
        }
        // 如果没有找到"flag=clash",添加该参数并更新URL
        if !clash_tag{
            params.Add("flag","clash")
            parsed_url.RawQuery = params.Encode()
            links[i].url = parsed_url.String()
        }
    }
    // 返回处理后的URL列表和nil错误
    return links,nil
}
// config_merge 根据模板和是否合并所有配置的标志,来合并配置
// template: 配置模板的字符串表示
// all: 是否合并所有配置的布尔值
// 返回错误信息,如果合并过程中发生错误
func config_merge(template string,all bool) error{
    // 从模板中提取日志、DNS、入站和实验性配置
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
    // 创建一个新的JSON对象,用于存储合并后的配置
    config := simplejson.New()
    config.Set("log", log)
    config.Set("dns", dns)
    config.Set("inbounds", inbounds)
    config.Set("experimental", experimental)
    // 格式化URL,并为每个URL配置创建一个错误通道
    links,err := format_url()
    error_channel := make(chan error,len(links))
    var jobs sync.WaitGroup
    if err != nil{
        return err
    }
    // 并发处理每个URL链接的配置合并
    for i,link := range links{
        // 如果不合并所有配置,只处理最后一个URL
        if !all && i != len(links)-1{
            continue
        }
        jobs.Add(1)
        go func(link config_link,template string,config *simplejson.Json,index int) {
            defer jobs.Done()
            // 克隆配置对象,以确保每个URL配置的独立性
            full_config := clone.Clone(config)
            project_dir,err := utils.Get_value("project-dir")
            if err != nil{
                utils.Logger_caller("Get project dir failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url config failed",index)
                return
            }
            // 合并路由配置
            route,err := Merge_route(template,link.url,link.proxy)
            if err != nil{
                utils.Logger_caller("Get route failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url failed,config:%s",index,link.label)
                return
            }
            full_config.(*simplejson.Json).Set("route", route)
            // 合并出站配置
            proies,err := Merge_outbounds(link.url,template)
            if err != nil{
                utils.Logger_caller("Get outbounds failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url failed,config:%s",index,link.label)
                return
            }
            full_config.(*simplejson.Json).Set("outbound", proies)
            // 对合并后的配置进行编码
            config_bytes,_ := full_config.(*simplejson.Json).EncodePretty()
            // 对标签进行MD5加密
            label,err := encryption_md5(link.label)
            if err != nil{
                utils.Logger_caller("Encryption md5 failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url failed,config:%s",index,link.label)
                return
            }
            // 将配置写入文件
            if err = utils.File_write(config_bytes,filepath.Join(project_dir.(string),"static",template,fmt.Sprintf("%s.json",label)),[]fs.FileMode{0666,0777});err != nil{
                utils.Logger_caller("Write config file failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url failed,config:%s",index,link.label)
                return
            }
        }(link,template,config,i)
    }
    // 等待所有并发任务完成
    jobs.Wait()
    close(error_channel)
    // 处理并输出任何配置合并过程中发生的错误
    for err := range error_channel{
        fmt.Println(err)
    }
    return nil
}
func Config_workflow(template string,all bool) error {
	if err := utils.Load_template(template); err != nil {
		utils.Logger_caller("load the template failed",err,1)
		return fmt.Errorf("load the %s template failed",template)
	}
	if err := utils.Load_config("Proxy"); err != nil {
		utils.Logger_caller("load the Proxy config failed",err,1)
		return fmt.Errorf("load the Proxy config failed")
	}
	config_merge(template,all)
	return nil
}
