package singbox

import (
	"fmt"
	"net/url"
	"sifu-box/models"

	"go.uber.org/zap"
)

// formatProviderURL 函数处理一个Provider切片, 为每个远程机场的URL添加"clash"标签参数(如果尚未存在)
// 这个函数接收一个Provider切片和一个Logger对象, 用于记录处理过程中的错误信息
// 它返回一个更新后的Provider切片和一个错误切片, 包含处理过程中遇到的任何错误
func formatProviderURL(providers []models.Provider, logger *zap.Logger) ([]models.Provider, []error) {
    // 初始化一个错误切片, 用于存储处理过程中遇到的错误
    var errors []error

    // 遍历Provider切片
    for i, provider := range providers {
        // 检查Provider是否为远程类型
        if provider.Remote {
            // 尝试解析Provider的路径URL
            providerURL, err := url.Parse(provider.Path)
            if err != nil {
                // 如果解析失败,记录错误日志并添加错误到错误切片
                logger.Error(fmt.Sprintf("解析'%s'链接失败: [%s]", provider.Name, err.Error()))
                errors = append(errors, fmt.Errorf("解析'%s'链接失败", provider.Name))
                continue
            }

            // 获取URL查询参数
            params := providerURL.Query()
            // 初始化一个标志变量,用于标记是否已存在"clash"标签
            clashTag := false

            // 遍历查询参数,检查是否存在"flag"为"clash"的参数
            for key, values := range params {
                if key == "flag" && values[0] == "clash" {
                    clashTag = true
                    break
                }
            }

            // 如果不存在"clash"标签,则添加该标签到查询参数中
            if !clashTag {
                params.Add("flag", "clash")
                // 更新URL的查询参数部分
                providerURL.RawQuery = params.Encode()
                // 更新Provider的路径为带有“clash”标签的URL
                providers[i].Path = providerURL.String()
            }
        }
    }
    // 返回更新后的Provider切片和错误切片
    return providers, errors
}

