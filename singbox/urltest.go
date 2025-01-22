package singbox

import (
	"sifu-box/models"

	"go.uber.org/zap"
)

// addURLTestOutbound 创建并添加一个URL测试出站配置到出站列表中
// 此函数接收一个出站配置列表和一个标签列表, 以及一个日志记录器
// 它会生成一个新的URL测试出站配置, 将其添加到出站配置列表中, 并在标签列表中添加相应的标签
// 如果生成过程中发生错误, 函数将返回该错误
// 参数:
//   outbounds - 当前的出站配置列表
//   tags - 当前的标签列表
//   logger - 用于记录日志的记录器
// 返回值:
//   []models.Outbound - 更新后的出站配置列表
//   []string - 更新后的标签列表
//   error - 如果操作成功, 则返回nil; 否则返回发生的错误
func addURLTestOutbound(outbounds []models.Outbound, tags []string, logger *zap.Logger) ([]models.Outbound, []string, error){
    // 初始化一个URL测试出站配置对象
    var urlTest models.URLTest
    // 定义URL测试出站配置的映射, 包含类型、是否中断已存在的连接、标签和出站标签列表
    URLTestMap := map[string]interface{}{"type":"urltest", "interrupt_exist_connections":false, "tag":"auto", "outbounds": tags}
    // 将URL测试出站配置对象视为一个一般的出站配置对象
    var outbound models.Outbound = &urlTest
    // 调用Transform方法, 根据提供的映射和日志记录器, 将出站配置对象转换为最终形式
    outbound, err := outbound.Transform(URLTestMap, logger)
    // 如果转换过程中发生错误, 返回错误
    if err != nil {
        return nil, nil, err
    } 
    // 将新出站配置的标签添加到标签列表中
    tags = append(tags, outbound.GetTag())
    // 将新的出站配置添加到出站配置列表中
    outbounds = append(outbounds, outbound)
    // 返回更新后的出站配置列表、标签列表和nil错误
    return outbounds, tags, nil
}