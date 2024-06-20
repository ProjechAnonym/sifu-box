package controller

import (
	"fmt"
	"net"
	"net/url"
	utils "sifu-box/Utils"
)

// Is_localhost 检查输入的URL是否指向本地主机
// 它返回一个布尔值表示是否为本地主机,以及可能的错误
func Is_localhost(input_url string) (bool,error) {
    // 解析输入的URL
	parsed_url,err := url.Parse(input_url)
	if err != nil {
        // 日志记录URL解析错误
		utils.Logger_caller("Error parsing URL",err,1)
		return false,fmt.Errorf("error parsing url")
	}
    // 获取URL的主机名
	host := parsed_url.Hostname()
    // 尝试将主机名解析为IP地址
	if ip := net.ParseIP(host); ip != nil {
        // 检查是否为回环地址,回环地址不被允许
		if ip.IsLoopback() {
			utils.Logger_caller("Deny!",fmt.Errorf("loopback address is not allowed"),1)
			return false, fmt.Errorf("loopback address is not allowed")
		}
        // 检查是否为IPv4地址
		if ip.To4() != nil {
            // 获取本地网络接口地址
			ips,err := net.InterfaceAddrs()
			if err != nil {
                // 日志记录获取接口地址失败
				utils.Logger_caller("Error",fmt.Errorf("get interface failed"),1)
				return false, fmt.Errorf("get interface failed")
			}
            // 遍历所有接口地址,检查是否与输入的IP地址匹配
			for _,addr := range ips{
				ip_addr,_,err := net.ParseCIDR(addr.String())
				if err != nil {
                    // 日志记录CIDR解析失败
					utils.Logger_caller("Error",fmt.Errorf("parse cidr failed"),1)
					return false, fmt.Errorf("parse cidr failed")
				}
				// 匹配本机ip,确认输入的url指向本机
				if ip.Equal(ip_addr) {
					return true, nil
				}
			}
			return false,nil
		}
        // 如果不是IPv4地址,则返回错误,不支持IPv6地址
		return false, fmt.Errorf("ipv6 address is not allowed")
	} 
    // 如果不是IP地址,则返回错误,域名不被允许。
	return false, fmt.Errorf("domain is not allowed")
}