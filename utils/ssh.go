package utils

import (
	"fmt"
	"net"
	"net/url"
	"sifu-box/models"

	"golang.org/x/crypto/ssh"
)

// tempKeyCallback is a callback function for adding a temporary key to the known_hosts file.
// Its main purpose is to update or verify the fingerprint of a host in the database based on the provided host and key.
// Parameters:
// host: Information about the host to which the temporary key belongs.
// key: The public key of the temporary key.
// Return value:
// If the operation is successful, it returns nil; otherwise, it returns the corresponding error.
func tempKeyCallback(host models.Host, key ssh.PublicKey) error {
    // Check if the host's fingerprint is empty, if so, update the fingerprint.
    if host.Fingerprint == "" {
        // Calculate the fingerprint of the provided key using SHA256.
        fingerPrint := ssh.FingerprintSHA256(key)
        
        // Update the host's fingerprint in the database. If the update fails, return the error.
        if err := DiskDb.Model(&host).Where("url = ?",host.Url).Update("fingerprint",fingerPrint).Error; err != nil {
            return err
        }
    } else {
        // If the host already has a fingerprint, compare it with the fingerprint of the provided key.
        if host.Fingerprint != ssh.FingerprintSHA256(key) {
            // If the fingerprints do not match, return an error message indicating a fingerprint mismatch.
            return fmt.Errorf("fingerprint mismatch")
        }
    }
    
    // If the process completes successfully, return nil.
    return nil
}

// InitClient 初始化SSH客户端配置。
// 参数host为待连接的服务器信息。
// 返回值为ssh.ClientConfig类型的指针，用于SSH客户端初始化，
// 返回值为服务器地址和端口的组合字符串，
// 返回值为错误信息，如果初始化过程中出现错误。
func InitClient(host models.Host) (*ssh.ClientConfig,string,error) {
    
    // 解析服务器URL。
    hostUrl,err := url.Parse(host.Url)
    if err != nil{
        // 如果URL解析失败，返回错误。
        return nil,"",err
    }
    
    // 构造SSH服务地址，使用标准SSH端口22。
    addr := hostUrl.Hostname() + ":22"
    
    // 创建SSH客户端配置对象。
    config := &ssh.ClientConfig{
        // 设置SSH登录用户名。
        User: host.Username,
        // 设置SSH认证方法，使用密码认证。
        Auth: []ssh.AuthMethod{ssh.Password(host.Password)},
        
        // 设置主机密钥回调验证函数。
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            
            // 使用临时的密钥回调函数验证主机密钥。
            if err := tempKeyCallback(host,key); err != nil {
                // 如果验证失败，返回错误。
                return err
            }
            // 验证成功，返回nil。
            return nil
        },
    }
    
    // 返回SSH客户端配置、地址和端口组合以及可能的错误。
    return config,addr,nil
}