package controller

import (
	"fmt"
	"sifu-box/models"
	"sifu-box/utils"
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