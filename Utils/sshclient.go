package utils

import (
	"fmt"
	"net"
	"net/url"

	"golang.org/x/crypto/ssh"
)

// temp_key_callback 是一个回调函数,用于在SSH连接过程中验证服务器的公钥
// 它通过比较预期的指纹和实际的公钥指纹来确保连接的安全性
// 如果服务器的指纹尚未记录,则会将其更新到数据库中
//
// 参数:
// - server: 代表服务器的信息,包括URL和当前的指纹
// - key: SSH连接过程中客户端提供的公钥
//
// 返回值:
// - 如果公钥验证失败或更新数据库时出错,则返回一个错误；否则返回nil
func temp_key_callback(server Server, key ssh.PublicKey) error {
    // 检查服务器是否有记录的指纹信息
	if server.Fingerprint == ""{
        // 计算并获取公钥的SHA256指纹
		finger_print := ssh.FingerprintSHA256(key)
        // 更新数据库中服务器的指纹信息
		if err := Db.Model(&server).Where("url = ?",server.Url).Update("fingerprint",finger_print).Error; err != nil{
            // 记录更新指纹信息失败的日志
			Logger_caller("update fingerprint failed",err,1)
			return err
		}
	}else{
        // 如果服务器有记录的指纹信息,则与提供的公钥指纹进行比较
		if server.Fingerprint != ssh.FingerprintSHA256(key){
            // 如果指纹不匹配,则记录错误日志并返回错误
			Logger_caller("fingerprint mismatch",nil,1)
			return fmt.Errorf("fingerprint mismatch")
		}
	}
    // 公钥验证成功,返回nil
	return nil
}
// Init_sshclient 初始化SSH客户端配置
// 该函数根据服务器信息解析URL,配置SSH连接参数,并返回SSH客户端配置及服务器地址
// 参数server 包含服务器的URL、用户名和密码
// 返回值为SSH客户端配置、服务器地址以及可能的错误
func Init_sshclient(server Server) (*ssh.ClientConfig,string,error) {
    // 解析服务器的URL以获取主机名
    host,err := url.Parse(server.Url)
    if err != nil{
        // 日志记录URL解析失败的错误
        Logger_caller("parse server url failde",err,1)
        return nil,"",err
    }
    // 构造SSH服务器地址,默认端口为22
    addr := host.Hostname() + ":22"
    // 配置SSH客户端配置
    config := &ssh.ClientConfig{
        User: server.Username,
        Auth: []ssh.AuthMethod{ssh.Password(server.Password)},
        // 定制化HostKeyCallback,用于验证服务器的公钥
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            // 调用自定义的公钥验证函数
            if err := temp_key_callback(server,key); err != nil {
                return err
            }
            return nil
        },
    }
    // 返回配置好的SSH客户端配置、服务器地址和nil错误
    return config,addr,nil
}