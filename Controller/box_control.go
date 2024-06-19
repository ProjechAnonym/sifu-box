package controller

import (
	"fmt"
	"io/fs"
	"path/filepath"
	singbox "sifu-box/Singbox"
	utils "sifu-box/Utils"

	"gopkg.in/yaml.v3"
)

// Add_items 向配置文件中添加新项目或更新现有项目配置
// box_config: 包含待添加或更新的项目配置的数据结构
// 返回值: 如果操作成功,则返回nil；否则返回错误信息
func Add_items(box_config utils.Box_config) error {
    // 获取项目目录路径,用于确定生成文件的路径
    project_dir, err := utils.Get_value("project-dir")
    if err != nil {
        // 记录获取项目目录失败的日志并返回错误
        utils.Logger_caller("Get project dir failed", err, 1)
        return fmt.Errorf("get project dir failed")
    }

    // 加载代理配置
    if err := utils.Load_config("Proxy"); err != nil {
        // 记录加载代理配置失败的日志并返回错误
        utils.Logger_caller("Load proxy config failed", err, 1)
        return fmt.Errorf("load proxy failed")
    }

    // 从配置中获取代理设置
    proxy_config, err := utils.Get_value("Proxy")
    proxyConfig := proxy_config.(utils.Box_config)
	// 初始化变量
	// index 用于存储新url的索引,为空则更新所有url的链接
	var index []int
	// urls 用于存储原url并添加新url
    urls := proxyConfig.Url
	// rulesets 用于存储原规则集并添加新规则集
    rulesets := proxyConfig.Rule_set
	
    if err != nil {
        // 记录获取代理配置失败的日志并返回错误
        utils.Logger_caller("Get proxy config failed", err, 1)
        return fmt.Errorf("get Proxy failed")
    }

    // 根据新添加的规则集和URL更新配置
	// 规则集为空则只更新新的url配置
    if len(box_config.Rule_set) == 0 {
		// 确定原先url长度,方便接下来进行指定url的更新
        urls_length := len(urls)
		// url为空说明没有添加,返回错误
        if len(box_config.Url) == 0 {
            return fmt.Errorf("no new links")
        } else {
            // 添加新URL到列表,并更新要刷新配置的url索引
            urls = append(urls, box_config.Url...)
            for i := range box_config.Url {
                index = append(index, urls_length+i)
            }
        }
    } else {
        // 如果指定了规则集,将其与代理配置中的规则集合并
        rulesets = append(rulesets, box_config.Rule_set...)
        if len(box_config.Url) != 0 {
            // 如果指定了URL,将其与代理配置中的URL合并
            urls = append(urls, box_config.Url...)
        }
    }

    // 创建新的代理配置
    var new_proxy_config utils.Box_config
    new_proxy_config.Rule_set = rulesets
    new_proxy_config.Url = urls

    // 将新的代理配置转换为YAML格式
    new_proxy_yaml, err := yaml.Marshal(new_proxy_config)
    if err != nil {
        // 记录转换配置失败的日志并返回错误
        utils.Logger_caller("Marshal proxy config failed", err, 1)
        return fmt.Errorf("marshal Proxy failed")
    }

    // 更新代理配置文件
    if err := utils.File_write(new_proxy_yaml, filepath.Join(project_dir.(string), "config", "Proxy.config.yaml"), []fs.FileMode{0644, 0644}); err != nil {
        // 记录写入配置文件失败的日志并返回错误
        utils.Logger_caller("Write Proxy config failed!", err, 1)
        return err
    }

    // 配置工作流刷新配置文件
    if err := singbox.Config_workflow(index); err != nil {
        // 记录配置工作流失败的日志并返回错误
        utils.Logger_caller("Config workflow failed", err, 1)
        return fmt.Errorf("config workflow failed")
    }

    // 操作成功,返回nil
    return nil
}