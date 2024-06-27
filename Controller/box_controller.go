package controller

import (
	"fmt"
	"io/fs"
	"path/filepath"
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

    // 从配置中获取代理设置
    proxy_config, err := utils.Get_value("Proxy")
	// 初始化变量
	// urls 用于存储原url并添加新url
    urls := proxy_config.(utils.Box_config).Url
	// rulesets 用于存储原规则集并添加新规则集
    rulesets := proxy_config.(utils.Box_config).Rule_set
	
    if err != nil {
        // 记录获取代理配置失败的日志并返回错误
        utils.Logger_caller("Get proxy config failed", err, 1)
        return fmt.Errorf("get Proxy failed")
    }
    // 如果两个均为空则说明没有配置
    if len(box_config.Url) == 0 && len(box_config.Rule_set) == 0{
        return fmt.Errorf("no config")
    }
    // 遍历新添加的链接,确定哪个是可以添加的
    for _,box_url := range box_config.Url{
        // 添加标志,默认不添加
        add_tag := true
        for _,url := range proxy_config.(utils.Box_config).Url{
            // 遍历已有的链接,如果标签一致则不添加
            if box_url.Label == url.Label{
                add_tag = false
                break
            }
        }
        if add_tag{
            urls = append(urls, box_url)
        }
    }
    // 遍历新添加的规则集,确定哪个是可以添加的
    for _,box_ruleset := range box_config.Rule_set{
        // 添加标志,默认不添加
        add_tag := true
        for _,ruleset := range proxy_config.(utils.Box_config).Rule_set{
            // 遍历已有的链接,如果标签一致则不添加
            if box_ruleset.Label == ruleset.Label{
                add_tag = false
                break
            }
        }
        if add_tag{
            rulesets = append(rulesets, box_ruleset)
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
    // 重新设置代理配置
    if err := utils.Set_value(new_proxy_config,"Proxy"); err != nil {
        // 记录设置代理配置失败的日志并返回错误
        utils.Logger_caller("Set proxy config failed", err, 1)
        return fmt.Errorf("set proxy failed")
    }
    // 操作成功,返回nil
    return nil
}

// Fetch_items 从配置中加载并返回代理设置
// 如果加载或获取配置过程中出现错误,将返回错误信息
// 返回值:
//   utils.Box_config: 代理配置的结构体
//   error: 加载或获取配置时可能出现的错误
func Fetch_items() (utils.Box_config, error) {
    // 获取名为"Proxy"的配置项的值如果获取失败,记录错误日志并返回错误
    proxy_config, err := utils.Get_value("Proxy")
    if err != nil {
        utils.Logger_caller("Get proxy config failed", err, 1)
        return utils.Box_config{}, fmt.Errorf("get Proxy failed")
    }
    // 将获取到的配置值断言为utils.Box_config类型,并返回该配置及nil错误
    return proxy_config.(utils.Box_config), nil
}
// Delete_items 根据提供的items删除配置文件中的特定URL和规则集。
// items: 一个映射，包含要删除的URL和规则集的索引。
// 返回值: 删除操作可能返回的任何错误。
func Delete_items(items map[string][]int) error{
    // 获取项目目录路径,用于确定生成文件的路径
    project_dir, err := utils.Get_value("project-dir")
    if err != nil {
        // 记录获取项目目录失败的日志并返回错误
        utils.Logger_caller("Get project dir failed", err, 1)
        return fmt.Errorf("get project dir failed")
    } 
    // 从配置中获取代理设置
    proxy_config, err := utils.Get_value("Proxy")
    if err != nil {
        // 记录日志并返回错误，如果无法获取代理配置
        utils.Logger_caller("Get proxy config failed", err, 1)
        return fmt.Errorf("get Proxy failed")
    }
    
    // 检查要删除的URL数量是否合法
    // urls 用于存储原url并添加新url
    urls := proxy_config.(utils.Box_config).Url
    var status bool
    status = Check_array(items["urls"],len(urls))
    if !status {
        // 记录日志并返回错误，如果URL数量不合法
        utils.Logger_caller("length error",fmt.Errorf("length is too long or index is too big"),1)
        return fmt.Errorf("length error")
    }

    // 检查要删除的规则集数量是否合法
    // rulesets 用于存储原规则集并添加新规则集
    rulesets := proxy_config.(utils.Box_config).Rule_set
    status = Check_array(items["rulesets"],len(rulesets))
    if !status {
        // 记录日志并返回错误，如果规则集数量不合法
        utils.Logger_caller("length error",fmt.Errorf("length is too long or index is too big"),1)
        return fmt.Errorf("length error")
    }
    // 从规则集中删除指定的规则集
    new_rulesets := Remove_item(rulesets,items["rulesets"]).([]utils.Box_ruleset)
    // 从URL列表中删除指定的URL
    new_urls := Remove_item(urls,items["urls"]).([]utils.Box_url)
    // 删除配置文件中的旧URL和规则集
    if err := Delete_config(urls,items["urls"],project_dir.(string));err != nil{
        // 记录日志并返回错误，如果删除配置失败
        utils.Logger_caller("Delete config failed",err,1)
        return fmt.Errorf("delete config failed")
    }

    // 创建新的代理配置
    var new_proxy_config utils.Box_config
    new_proxy_config.Rule_set = new_rulesets
    new_proxy_config.Url = new_urls

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
    // 重新设置代理配置
    if err := utils.Set_value(new_proxy_config,"Proxy"); err != nil {
        // 记录设置代理配置失败的日志并返回错误
        utils.Logger_caller("Set proxy config failed", err, 1)
        return fmt.Errorf("set proxy failed")
    }

    if len(items["urls"]) != 0 {
        var servers []utils.Server
        if err := utils.Db.Model(&utils.Server{}).Find(&servers).Error;err != nil {
            // 记录查询服务器列表失败的日志并返回错误
            utils.Logger_caller("Query server list failed", err, 1)
            return fmt.Errorf("query server list failed")
        }
        for _,server := range(servers){
            change_tag := true
            for _,url_label := range(new_proxy_config.Url){
                if server.Config == url_label.Label{
                    change_tag = false
                    break
                }
            }
            if change_tag{
                if err := utils.Db.Model(&utils.Server{}).Where("url = ?",server.Url).Update("config",new_proxy_config.Url[0].Label).Error;err != nil{
                    return err
                }
            }
        }
    }
    // 删除操作成功，返回nil
    return nil
}