package controller

import (
	"fmt"
	"net"
	"net/url"
	"sifu-box/utils"
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