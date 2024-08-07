package utils

import (
	"io"
	"os"
	"path/filepath"
	"sifu-box/models"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SftpRead 通过SFTP协议从远程主机读取文件内容
// 参数host表示远程服务器的连接信息,src为要读取的文件路径
// 返回文件内容的字节切片和可能的错误
func SftpRead(host models.Host, src string) ([]byte, error) {
    // 初始化客户端配置、地址和错误处理
    config, addr, err := InitClient(host)
    if err != nil {
        // 如果初始化失败,返回错误
        return nil, err
    }

    // 基于SSH协议建立与远程主机的连接
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        // 如果连接失败,返回错误
        return nil, err
    }
    // 确保在函数返回前关闭SSH客户端连接
    defer client.Close()

    // 创建SFTP客户端
    sftpClient, err := sftp.NewClient(client)
    if err != nil {
        // 如果创建SFTP客户端失败,返回错误
        return nil, err
    }
    // 确保在函数返回前关闭SFTP客户端连接
    defer sftpClient.Close()

    // 打开要读取的文件
    srcFile, err := sftpClient.OpenFile(src, os.O_RDONLY)
    if err != nil {
        // 如果文件打开失败,返回错误
        return nil, err
    }
    // 确保在函数返回前关闭文件
    defer srcFile.Close()

    // 读取文件全部内容
    content, err := io.ReadAll(srcFile)
    if err != nil {
        // 如果读取文件内容失败,返回错误
        return nil, err
    }
    
    // 返回文件内容和无错误信息
    return content, nil
}

// SftpWrite 通过SFTP协议将内容写入到远程主机上的指定文件
// host: 远程主机的模型,包含连接必要的信息
// content: 需要写入的字节切片内容
// dst: 远程主机上的目标文件路径
// 返回错误信息,如果操作成功,则为nil
func SftpWrite(host models.Host, content []byte, dst string) error {
    // 初始化客户端配置、地址和错误处理
    config, addr, err := InitClient(host)
    if err != nil {
        return err
    }

    // 建立SSH客户端连接
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        return err
    }
    defer client.Close()

    // 创建SFTP客户端
    sftpClient, err := sftp.NewClient(client)
    if err != nil {
        return err
    }
    defer sftpClient.Close()

    // 检查目标文件的目录是否存在,如果不存在则创建
    if _, err := sftpClient.Stat(filepath.Dir(dst)); err != nil {
        if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile) {
            if err := sftpClient.MkdirAll(filepath.Dir(dst)); err != nil {
                return err
            }
        } else {
            return err
        }
    }

    // 打开目标文件,如果不存在则创建,同时支持读写和截断操作
    file, err := sftpClient.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
    defer func() {
        if err := file.Close(); err != nil {
            LoggerCaller("文件无法关闭", err, 1)
        }
    }()
    if err != nil {
        return err
    }

    // 写入内容到文件
    _, err = file.Write(content)
    if err != nil {
        return err
    }

    // 返回nil表示文件写入成功
    return nil
}

// SftpDelete 通过 SFTP 协议删除指定主机上的文件或目录
// host: 包含主机信息的模型,包括IP、端口等
// path: 需要删除的文件或目录的路径
// 返回错误信息,如果操作成功,则返回nil
func SftpDelete(host models.Host, path string) error {
    // 初始化SSH客户端配置、地址和错误处理
    config, addr, err := InitClient(host)
    if err != nil {
        return err
    }

    // 基于配置和地址,通过SSH协议建立连接
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        return err
    }
    // 确保在函数返回前关闭客户端连接
    defer client.Close()

    // 创建SFTP客户端以进行文件操作
    sftpClient, err := sftp.NewClient(client)
    if err != nil {
        return err
    }
    // 确保在函数返回前关闭SFTP客户端连接
    defer sftpClient.Close()

    // 检查指定路径的文件或目录是否存在
    if _, err := sftpClient.Stat(path); err != nil {
        // 如果是文件或目录不存在的错误,则直接返回成功
        if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile) {
            return nil
        } else {
            // 其他类型的错误进行返回
            return err
        }
    }

    // 删除指定路径的文件或目录,包括其内容
    if err := sftpClient.RemoveAll(path); err != nil {
        return err
    }

    // 操作成功,返回nil
    return nil
}