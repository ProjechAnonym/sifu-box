package controller

import (
	"fmt"
	"net"
	"net/url"
	"sifu-box/execute"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"
)

// IsLocalhost 检查给定的URL是否指向本地主机
// 参数:
//   input_url - 需要检查的URL字符串
// 返回值:
//   bool - 如果URL指向本地主机则返回true,否则返回false
//   error - 如果检查过程中发生错误,返回相应的错误信息
func IsLocalhost(input_url string) (bool, error) {
	
	// 解析输入的URL
	parsedUrl, err := url.Parse(input_url)
	if err != nil {
		// 如果URL解析失败,记录错误并返回
		utils.LoggerCaller("无法解析url", err, 1)
		return false, fmt.Errorf("无法解析url")
	}
	
	// 获取URL的主机部分
	host := parsedUrl.Hostname()
	
	// 解析主机为IP地址
	if ip := net.ParseIP(host); ip != nil {
		
		// 检查是否为回环地址
		if ip.IsLoopback() {
			// 如果是回环地址,记录错误并返回
			utils.LoggerCaller("地址类型错误", fmt.Errorf("不允许设置回环地址"), 1)
			return false, fmt.Errorf("不允许设置回环地址")
		}
		
		// 检查是否为IPv4地址
		if ip.To4() != nil {
			
			// 获取网络接口地址
			ips, err := net.InterfaceAddrs()
			if err != nil {
				// 如果获取失败,记录错误并返回
				utils.LoggerCaller("获取接口失败", err, 1)
				return false, fmt.Errorf("获取接口失败")
			}
			
			// 遍历所有接口地址,检查是否与输入的IP匹配
			for _, addr := range ips {
				ip_addr, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					// 如果地址解析失败,记录错误并返回
					utils.LoggerCaller("解析地址失败", err, 1)
					return false, fmt.Errorf("解析地址失败")
				}
				
				// 如果IP地址匹配,则返回true
				if ip.Equal(ip_addr) {
					return true, nil
				}
			}
			// 如果遍历完所有接口地址都没有匹配的,返回false
			return false, nil
		}
		
		// 如果是IPv6地址,返回错误
		return false, fmt.Errorf("不支持ipv6")
	}
	
	// 如果输入的是域名,返回错误
	return false, fmt.Errorf("不支持域名")
}

// SwitchTemplate 根据提供的模板和URL列表更新主机的配置。
// 该函数首先从数据库中查询与提供的URL列表匹配的主机以及所有的提供商。
// 然后，它检查每个主机的配置是否在提供商列表中，如果不在，则更新主机的配置为第一个提供商的名称。
// 最后，它使用并行处理的方式更新主机的配置，并返回可能发生的错误列表。
// 参数:
//   template - 用于更新的模板字符串，目前未使用。
//   urls - 包含需要更新的主机URL的切片。
//   lock - 用于并行处理时同步的互斥锁指针。
// 返回值:
//   []error - 包含执行过程中可能发生的错误的切片。
func SwitchTemplate(template string,urls []string,lock *sync.Mutex) []error {
    var hosts []models.Host
    // 从磁盘数据库查询与提供的URL列表匹配的主机。
    if err := utils.DiskDb.Table("hosts").Where("url IN (?)", urls).Find(&hosts).Error; err != nil {
        utils.LoggerCaller("数据库查询失败", err, 1)
        return []error{fmt.Errorf("数据库查询失败")}
    }
    var providers []models.Provider
    // 从内存数据库查询所有提供商。
    if err := utils.MemoryDb.Find(&providers).Error; err != nil {
        utils.LoggerCaller("数据库查询失败", err, 1)
        return []error{fmt.Errorf("数据库查询失败")}
    }

    serverUpdate := false

    // 遍历所有查询到的主机。
    for _, host := range hosts {
        // 遍历所有提供商，检查主机的配置是否在提供商列表中。
        for _, provider := range providers {
            if host.Config == provider.Name {
                serverUpdate = true
                break
            }
        }

        // 如果主机的配置不在提供商列表中，则将其更新为第一个提供商的名称。
        if !serverUpdate {
            host.Config = providers[0].Name
            // 更新数据库中相应主机的配置。
            if err := utils.DiskDb.Model(&models.Host{}).Where("url = ?", host.Url).Update("config", providers[0].Name).Error; err != nil {
                utils.LoggerCaller("更新主机机场信息失败", err, 1)
                return []error{err}
            }
        }
    }
    // 使用并行处理的方式更新主机的配置。
    errs := execute.GroupUpdate(hosts, providers, lock, true)
    return errs
}