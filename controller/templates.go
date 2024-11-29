package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"sifu-box/models"
	"sifu-box/utils"

	"gopkg.in/yaml.v3"
)

// GetTemplates 获取所有模板配置
// 该函数从配置中获取模板信息，并将其转换为一个包含模板名称和信息的映射列表
// 返回值:
// - []map[string]interface{}: 包含模板名称和信息的映射列表
// - error: 错误信息，如果获取模板配置失败，则返回错误
func GetTemplates() ([]map[string]interface{}, error) {
    // 从配置中获取模板信息
    templatesMap, err := utils.GetValue("templates")
    if err != nil {
        // 如果获取模板配置失败，则记录错误信息并返回错误
        utils.LoggerCaller("获取模板配置失败", err, 1)
        return nil,fmt.Errorf("获取模板配置失败")
    }
    var templates []map[string]interface{}
    // 遍历模板信息，将其转换为包含模板名称和信息的映射，并添加到列表中
    for key, template := range templatesMap.(map[string]models.Template) {
        templates = append(templates, map[string]interface{}{"Name": key, "Template": template})
    }
    // 返回包含所有模板信息的列表
    return templates,nil
}

// AddTemplate 添加一个新的模板到项目中
// 参数:
//   name - 模板的名称
//   template - 模板的内容
// 返回值:
//   如果添加模板过程中发生错误，则返回错误
func AddTemplate(name string, template models.Template) error {
    // 获取项目目录路径
    projectDir,err := utils.GetValue("project-dir")
    if err != nil {
        utils.LoggerCaller("获取项目目录失败", err, 1)
        return fmt.Errorf("获取项目目录失败")
    }
    // 获取当前的模板配置
    templates, err := utils.GetValue("templates")
    if err != nil {
        utils.LoggerCaller("获取模板配置失败", err, 1)
        return fmt.Errorf("获取模板配置失败")
    }
    // 在模板配置中添加新的模板
    templates.(map[string]models.Template)[name] = template
    // 更新模板配置
    if err := utils.SetValue(templates, "templates");err != nil {
        utils.LoggerCaller("更新模板配置失败", err, 1)
        return fmt.Errorf("更新模板配置失败")
    }
    // 序列化模板为yaml格式
    templateYaml, err := yaml.Marshal(template)
    if err != nil {  
        utils.LoggerCaller("序列化yaml文件失败", err, 1)
        return fmt.Errorf("序列化yaml文件失败")
    }
    // 将模板文件写入到项目目录中
    if err := utils.FileWrite(templateYaml, filepath.Join(projectDir.(string), "template", fmt.Sprintf("%s.template.yaml",name))); err != nil { 
        utils.LoggerCaller("写入模板文件失败", err, 1)
        return fmt.Errorf("写入模板文件失败")
    }
    return nil
}

// RefreshTemplates 用于刷新模板，通过从项目目录中读取和解析模板文件来更新模板信息。
// 它主要关注于读取和解析名为 "recover.template.yaml" 的文件，该文件用于定义恢复操作的模板。
// 函数返回一个映射，其中键是模板名称，值是模板对象。如果在执行过程中遇到错误，将返回一个错误对象。
func RefreshTemplates() (map[string]models.Template,error) {
    // 尝试获取项目目录路径，这是后续读取配置文件的基础路径。
    projectDir,err := utils.GetValue("project-dir")
    if err != nil {
        // 如果获取项目目录失败，记录错误并返回提示信息。
        utils.LoggerCaller("获取项目目录失败", err, 1)
        return nil,fmt.Errorf("获取项目目录失败")
    }

    // 读取项目目录下的恢复模板文件。
    file,err := os.ReadFile(filepath.Join(projectDir.(string),"config","recover.template.yaml"))
    if err != nil {
        // 如果读取文件失败，记录错误并返回提示信息。
        utils.LoggerCaller("读取恢复模板文件失败", err, 1)
        return nil,fmt.Errorf("读取恢复模板文件失败")
    }

    // 解析读取到的文件内容到 Template 结构体。
    var recoverTemplate models.Template
    if err := yaml.Unmarshal(file,&recoverTemplate);err != nil {
        // 如果解析文件失败，记录错误并返回提示信息。
        utils.LoggerCaller("解析恢复模板文件失败", err, 1)
        return nil,fmt.Errorf("解析恢复模板文件失败")
    }

    // 返回解析成功的模板信息。
    return map[string]models.Template{"recover":recoverTemplate},nil
}

// DeleteTemplate 删除指定名称的模板文件和配置
// 参数:
//   names []string - 待删除模板的名称列表
// 返回值:
//   []error - 删除过程中可能发生的错误列表
func DeleteTemplate(names []string) []error{
    // 获取项目目录路径
    projectDir,err := utils.GetValue("project-dir")
    if err != nil {
        utils.LoggerCaller("获取项目目录失败", err, 1)
        return []error{fmt.Errorf("获取项目目录失败")}
    }

    // 获取模板配置
    templates, err := utils.GetValue("templates")
    if err != nil {
        utils.LoggerCaller("获取模板配置失败", err, 1)
        return []error{fmt.Errorf("获取模板配置失败")}
    }
    
    // 存储删除过程中发生的错误
    var errs []error
    for _,name := range names {
        // 禁止删除默认模板
        if name == "default" {
            utils.LoggerCaller("删除默认模板失败", fmt.Errorf("禁止删除默认模板"), 1)
            errs = append(errs, fmt.Errorf("禁止删除默认模板"))
            continue
        }
        // 删除模板文件
        if err := utils.FileDelete(filepath.Join(projectDir.(string), "template", fmt.Sprintf("%s.template.yaml", name))); err != nil{
            utils.LoggerCaller("删除模板文件失败", err, 1)
            errs = append(errs, fmt.Errorf("删除模板文件失败"))
        }
        // 删除此模板生成的配置文件
        if err := utils.FileDelete(filepath.Join(projectDir.(string), "static", name)); err != nil{
            utils.LoggerCaller("删除模板文件失败", err, 1)
            errs = append(errs, fmt.Errorf("删除模板文件失败"))
        }
        // 从模板配置中移除该模板
        delete(templates.(map[string]models.Template), name)
    }

    // 更新模板配置
    if err := utils.SetValue(templates, "templates");err != nil {
        utils.LoggerCaller("更新模板配置失败", err, 1)
        return []error{fmt.Errorf("更新模板配置失败")}
    }
    
    // 返回删除过程中发生的错误列表
    return errs
}