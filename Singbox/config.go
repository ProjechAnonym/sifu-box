package singbox

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	utils "sifu-box/Utils"
	"sort"
	"strings"
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
func format_url(index []int) ([]utils.Box_url,error) {

    // 从配置文件中获取URL列表
    proxy_config,err := utils.Get_value("Proxy")
    if err != nil || len(proxy_config.(utils.Box_config).Url) == 0{
        // 记录获取URL列表时的错误
        utils.Logger_caller("Get urls failed!",err,1)
        return nil,err
    }
    // 检查URL列表是否为空
    if len(proxy_config.(utils.Box_config).Url) == 0 {
        // 如果URL列表为空,创建并返回一个新的空列表和一个错误
        err = errors.New("get url list failed")
        utils.Logger_caller("Get urls failed!",err,1)
        return nil,err
    }
    // 初始化一个用于存储处理后URL的切片
    // 确认要获取的URL数量
    var links_length int
    sort.Ints(index)
    // 如果index长度为0则解析所有url
    if len(index) == 0{
        links_length = len(proxy_config.(utils.Box_config).Url)
        for i := range proxy_config.(utils.Box_config).Url{
            index = append(index, i)
        }
    // 长度合理则确定url切片长度
    }else if len(index) <= len(proxy_config.(utils.Box_config).Url) && index[len(index)-1] < len(proxy_config.(utils.Box_config).Url){
        links_length = len(index)
    }else{
        // 长度超过则报错
        utils.Logger_caller("error",fmt.Errorf("parsing more urls in the config"),1)
        return nil,err
    }
    links := make([]utils.Box_url,links_length)
    
    // 初始化一个标志,用于指示是否找到了带有"flag=clash"的URL
    clash_tag := false
    // 遍历URL列表,处理每个URL
    for i, value := range index {

		// 从URL映射中提取并设置代理标志和标签
		links[i].Proxy = proxy_config.(utils.Box_config).Url[value].Proxy
		links[i].Label = proxy_config.(utils.Box_config).Url[value].Label
        // 解析URL字符串
        parsed_url,err := url.Parse(proxy_config.(utils.Box_config).Url[value].Url)
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
                links[i].Url = parsed_url.String()
                clash_tag = true
                break
            }
        }
        // 如果没有找到"flag=clash",添加该参数并更新URL
        if !clash_tag{
            params.Add("flag","clash")
            parsed_url.RawQuery = params.Encode()
            links[i].Url = parsed_url.String()
        }
    }
    // 返回处理后的URL列表和nil错误
    return links,nil
}
// config_merge 根据模板和是否合并所有配置的标志,来合并配置
// template: 配置模板的字符串表示
// all: 是否合并所有配置的布尔值
func config_merge(template string,mode bool,index []int) {
    // 从模板中提取日志、DNS、入站和实验性配置
    log,err := utils.Get_value(template,"log")
    if err != nil{
        utils.Logger_caller("Get log failed!",err,1)
        return 
    }
    dns,err := utils.Get_value(template,"dns")
    if err != nil{
        utils.Logger_caller("Get dns failed!",err,1)
        return 
    }
    inbounds,err := utils.Get_value(template,"inbounds")
    if err != nil{
        utils.Logger_caller("Get inbounds failed!",err,1)
        return
    }
    experimental,err := utils.Get_value(template,"experimental")
    if err != nil{
        utils.Logger_caller("Get experimental failed!",err,1)
        return 
    }
    // 创建一个新的JSON对象,用于存储合并后的配置
    config := simplejson.New()
    config.Set("log", log)
    config.Set("dns", dns)
    config.Set("inbounds", inbounds)
    config.Set("experimental", experimental)
    // 格式化URL,并为每个URL配置创建一个错误通道
    links,err := format_url(index)
    if err != nil{
        return 
    }
    error_channel := make(chan error,len(links))
    // 创建进程组,避免程序过早退出
    var jobs sync.WaitGroup
    
    // 并发处理每个URL链接的配置合并
    for i,link := range links{
        jobs.Add(1)
        go func(link utils.Box_url,template string,config *simplejson.Json,index int,mode bool) {
            defer jobs.Done()
            // 克隆配置对象,以确保每个URL配置的独立性
            full_config := clone.Clone(config)
            project_dir,err := utils.Get_value("project-dir")
            if err != nil{
                utils.Logger_caller("Get project dir failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url of %s config failed",index,template)
                return
            }
            // 合并路由配置
            route,err := Merge_route(template,link.Url,link.Proxy)
            if err != nil{
                utils.Logger_caller("Get route failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url of %s failed,config:%s",index,template,link.Label)
                return
            }
            full_config.(*simplejson.Json).Set("route", route)
            // 合并出站配置
            proies,err := Merge_outbounds(link.Url,template)
            if err != nil{
                utils.Logger_caller("Get outbounds failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url of %s failed,config:%s",index,template,link.Label)
                return
            }
            full_config.(*simplejson.Json).Set("outbound", proies)
            // 对合并后的配置进行编码
            config_bytes,_ := full_config.(*simplejson.Json).EncodePretty()
            // 获取配置文件名
            var label string
            if mode{
                // 对标签进行MD5加密
                label,err = encryption_md5(link.Label)
                if err != nil{
                    utils.Logger_caller("Encryption md5 failed!",err,1)
                    error_channel <- fmt.Errorf("generate the %dth url of %s failed,config:%s",index,template,link.Label)
                    return
                }
            }else{
                label = link.Label
            }
            // 将配置写入文件
            if err = utils.File_write(config_bytes,filepath.Join(project_dir.(string),"static",template,fmt.Sprintf("%s.json",label)),[]fs.FileMode{0666,0777});err != nil{
                utils.Logger_caller("Write config file failed!",err,1)
                error_channel <- fmt.Errorf("generate the %dth url of %s failed,config:%s",index,template,link.Label)
                return
            }
            utils.Logger_caller(fmt.Sprintf("Generate the %s config of %s success!",link.Label,template),nil,1)
        }(link,template,config,i,mode)
    }
    // 等待所有并发任务完成
    jobs.Wait()
    close(error_channel)
    // 处理并输出任何配置合并过程中发生的错误
    for err := range error_channel{
        utils.Logger_caller("generate error",err,1)
    }
    
}
func Config_workflow(index []int) error {
    project_dir,err := utils.Get_value("project-dir")
    if err != nil{
        utils.Logger_caller("get project dir failed",err,1)
        return err
    }

	// 打开目录
	template_dir, err := os.Open(filepath.Join(project_dir.(string),"template"))
	if err != nil {
		utils.Logger_caller("failed to open template directory", err,1)
	}
	defer template_dir.Close()

	// 读取目录条目
	entries, err := template_dir.ReadDir(-1) // -1 表示读取所有条目
	if err != nil {
		utils.Logger_caller("failed to read template directory", err,1)
	}
    if err := utils.Load_config("Proxy"); err != nil {
		utils.Logger_caller("load the Proxy config failed",err,1)
		return fmt.Errorf("load the Proxy config failed")
	}
    // 创建进程组,避免程序过早退出
    var workflow sync.WaitGroup
    // 获取服务器配置
    server_config,err := utils.Get_value("Server")
    if err != nil{
        // 记录并返回可能出现的错误
        utils.Logger_caller("get server config failed",err,1)
        return fmt.Errorf("get server config failed")
    }
    // 获取运行模式
    server_mode := server_config.(utils.Server_config).Server_mode
	// 打印所有条目的名称
	for _, entry := range entries {
        template := strings.Split(entry.Name(), ".")[0]
        if err := utils.Load_template(template); err != nil {
            utils.Logger_caller("load the template failed",err,1)
            return fmt.Errorf("load the %s template failed",template)
        }
        workflow.Add(1)
        go func ()  {
            defer func ()  {
                workflow.Done()
                utils.Del_key(template)
            }()
            config_merge(template,server_mode,index)
        }()
	}
	workflow.Wait()
    utils.Del_key("Proxy")
	return nil
}
