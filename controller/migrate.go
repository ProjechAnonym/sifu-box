package controller

import (
	"fmt"
	"path/filepath"
	"sifu-box/models"
	"sifu-box/utils"

	"gopkg.in/yaml.v3"
)

// ExportInfo 导出系统信息。
// 该函数从数据库中获取主机、提供商、规则集和模板信息，并将其导出为YAML格式的字符串。
// 返回值:
// - string: 导出的信息的YAML字符串表示。
// - error: 如果获取数据或导出信息过程中发生错误，则返回错误。
func ExportInfo() (string,error){
    // 初始化主机列表
    var hosts []models.Host
    // 从数据库中获取主机信息
    if err := utils.DiskDb.Model(&hosts).Find(&hosts).Error; err != nil {
        // 如果获取主机信息失败，记录错误并返回错误
        utils.LoggerCaller("获取主机失败", err, 1)
        return "",fmt.Errorf("获取主机失败")
    }

    // 初始化机场列表
    var providers []models.Provider
    // 从数据库中获取机场信息
    if err := utils.MemoryDb.Find(&providers).Error; err != nil {
        // 如果获取机场信息失败，记录错误并返回错误
        utils.LoggerCaller("获取代理信息失败", err, 1)
        return "",fmt.Errorf("获取代理信息失败")
    }

    // 初始化规则集列表
    var rulesets []models.Ruleset
    // 从数据库中获取规则集信息
    if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
        // 如果获取规则集信息失败，记录错误并返回错误
        utils.LoggerCaller("获取代理信息失败", err, 1)
        return "",fmt.Errorf("获取代理信息失败")
    }

    // 从配置中获取模板信息
    templates,err := utils.GetValue("templates"); 
    if err != nil {
        // 如果获取模板信息失败，记录错误并返回错误
        utils.LoggerCaller("获取模板信息失败", err, 1)
        return "",fmt.Errorf("获取模板信息失败")
    }

    // 初始化迁移信息映射
    var migrateInfo = models.Migrate{Hosts: hosts,Providers: providers,Rulesets: rulesets,Templates: templates.(map[string]models.Template)}

    // 将信息映射转换为YAML格式的字节切片
    migrateBytes,err := yaml.Marshal(migrateInfo)
    if err != nil {
        // 如果导出信息失败，记录错误并返回错误
        utils.LoggerCaller("导出信息失败", err, 1)
        return "",fmt.Errorf("导出信息失败")
    }

    // 返回YAML格式的信息字符串
    return string(migrateBytes),nil
}

// ImportInfo 导入项目信息。
// 该函数接收一个 models.Migrate 类型的参数 info，用于迁移项目中的主机、机场、规则集和模板信息。
// 返回值为 error 类型，用于指示在迁移过程中可能发生的错误。
func ImportInfo(info models.Migrate) error {
    // 获取项目目录路径。
    projectDir, err := utils.GetValue("project-dir")
    if err != nil {
        utils.LoggerCaller("获取项目目录失败", err, 1)
        return fmt.Errorf("获取项目目录失败")
    }
    
    // 导入主机信息。
    hosts := info.Hosts
    if len(hosts) != 0 {
        if err := utils.DiskDb.Create(&hosts).Error; err != nil {
            utils.LoggerCaller("导入主机信息失败", err, 1)
            return fmt.Errorf("主机信息写入数据库失败")
        }
    }
    
    // 导入机场信息。
    providers := info.Providers
    if len(providers) != 0 {
        if err := utils.MemoryDb.Create(&providers).Error; err != nil {
            utils.LoggerCaller("导入机场信息失败", err, 1)
            return fmt.Errorf("机场信息写入数据库失败")
        }
    }
    
    // 导入规则集信息。
    rulesets := info.Rulesets
    if len(rulesets) != 0 {
        if err := utils.MemoryDb.Create(&rulesets).Error; err != nil {
            utils.LoggerCaller("导入规则集信息失败", err, 1)
            return fmt.Errorf("规则集信息写入数据库失败")
        }
    }
    
    // 验证机场信息是否成功导入。
    if err := utils.MemoryDb.Find(&providers).Error; err != nil {
        utils.LoggerCaller("导入规则集信息失败", err, 1)
        return fmt.Errorf("规则集信息写入数据库失败")
    }
    
    // 验证规则集信息是否成功导入。
    if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
        utils.LoggerCaller("导入规则集信息失败", err, 1)
        return fmt.Errorf("规则集信息写入数据库失败")
    }
    
    // 序列化并写入 proxy 配置文件。
    proxy := models.Proxy{Providers: providers, Rulesets: rulesets}
    proxyYaml, err := yaml.Marshal(proxy)
    if err != nil {  
        utils.LoggerCaller("序列化yaml文件失败", err, 1)
        return fmt.Errorf("序列化yaml文件失败")
    }
    if err := utils.FileWrite(proxyYaml, filepath.Join(projectDir.(string), "config", "proxy.config.yaml")); err != nil { 
        utils.LoggerCaller("写入proxy配置文件失败", err, 1)
        return fmt.Errorf("写入proxy配置文件失败")
    }
    
    // 获取并更新模板配置。
    templates,err := utils.GetValue("templates")
    if err != nil {
        utils.LoggerCaller("获取模板配置失败", err, 1)
        return fmt.Errorf("获取模板配置失败")
    }
    for key, template := range info.Templates {
        templates.(map[string]models.Template)[key] = template
        
        // 序列化并写入每个模板文件。
        templateYaml, err := yaml.Marshal(template)
        if err != nil {  
            utils.LoggerCaller("序列化yaml文件失败", err, 1)
            return fmt.Errorf("序列化yaml文件失败")
        }
        if err := utils.FileWrite(templateYaml, filepath.Join(projectDir.(string), "template", fmt.Sprintf("%s.template.yaml",key))); err != nil { 
            utils.LoggerCaller("写入模板文件失败", err, 1)
            return fmt.Errorf("写入模板文件失败")
        }
    }
    
    // 更新模板信息到配置中。
    if err := utils.SetValue(templates, "templates");err != nil {
        utils.LoggerCaller("写入模板信息失败", err, 1)
        return fmt.Errorf("写入模板信息失败")
    }
    
    return nil
}