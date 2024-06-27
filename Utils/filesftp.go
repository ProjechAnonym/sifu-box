package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Sftp_delete 通过SFTP协议删除指定服务器上的文件或目录
// server: 服务器配置信息
// path: 需要删除的文件或目录路径
// 返回值: 删除操作可能产生的错误
func Sftp_delete(server Server,path string) error{
    // 初始化SSH客户端配置
	config,addr,err := Init_sshclient(server)
	if err != nil {
		// 记录初始化SSH配置失败的日志
		Logger_caller("init ssh config failed",err,1)
		return err
	}

    // 建立SSH连接
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		// 记录连接失败的日志
		Logger_caller("Failed to dial: ", err,1)
		return err
	}
	defer client.Close()

    // 初始化SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		// 记录初始化SFTP客户端失败的日志
		Logger_caller("Failed to init sftp client: ", err,1)
		return err
	}
	defer sftpClient.Close()

    // 检查文件或目录是否存在
	if _,err := sftpClient.Stat(path);err != nil{
		// 如果文件或目录不存在,记录日志并返回nil
		if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile){
			Logger_caller("file not found",err,1)
			return nil
		}else{
			// 其他错误,记录日志并返回错误
			Logger_caller("stat file failed",err,1)
			return err
		}
	}

    // 删除文件或目录
	if err := sftpClient.RemoveAll(path);err != nil{
		// 记录删除失败的日志
		Logger_caller("remove file failed",err,1)
		return err
	}

    // 删除操作成功,返回nil
	return nil
}
// Sftp_write 通过SFTP协议向指定服务器的文件写入内容
// server: 服务器配置信息
// content: 需要写入文件的内容
// dst: 目标文件路径
// 返回值: 错误信息,如果操作成功则为nil
func Sftp_write(server Server,content []byte, dst string) error{
    // 初始化SSH客户端配置并获取服务器地址
    config,addr,err := Init_sshclient(server)
    if err != nil {
        // 记录初始化SSH客户端时的错误
        Logger_caller("init ssh config failed",err,1)
        return err
    }
    // 建立SSH连接
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        Logger_caller("Failed to dial: ", err,1)
		return err
    }
    defer client.Close()
    // 初始化SFTP客户端
    sftpClient, err := sftp.NewClient(client)
    if err != nil {
        Logger_caller("Failed to init sftp client: ", err,1)
		return err
    }
    defer sftpClient.Close()
    // 检查目标文件目录是否存在,不存在则创建
    if _,err := sftpClient.Stat(filepath.Dir(dst)); err != nil{
        if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile){
            if err := sftpClient.MkdirAll(filepath.Dir(dst));err != nil {
                Logger_caller("mkdir failed",err,1)
                return err
            }
        }else{
            Logger_caller("stat dir failed",err,1)
            return err
        }
    }
    // 打开目标文件,如果不存在则创建,并设置为可写模式
    file, err := sftpClient.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
    defer func() {
        // 确保文件关闭
        if err := file.Close(); err != nil {
            Logger_caller("File can not close!", err,1)
        }
    }()
    if err != nil {
        Logger_caller("Create file failed", err,1)
        return err
    }
    // 写入内容到文件
    _, err = file.Write(content)
    if err != nil {
        Logger_caller("Write config failed!", err,1)
        return err
    }

    // 操作成功
    return nil
}

// Sftp_read 通过SFTP协议从指定服务器读取文件内容
// server: 服务器配置信息,包括SSH连接所需的用户名、密码、主机等
// src: 指定文件的路径
// 返回读取到的文件内容以及可能发生的错误
func Sftp_read(server Server,src string) ([]byte,error){
    // 初始化SSH客户端配置
	config,addr,err := Init_sshclient(server)
	if err != nil {
		// 记录初始化SSH客户端配置失败的日志
		Logger_caller("init ssh config failed",err,1)
		return nil,err
	}

    // 建立SSH连接
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		// 记录连接服务器失败的日志
		Logger_caller("Failed to dial: ", err,1)
		return nil,err
	}
	defer client.Close()

    // 初始化SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		// 记录初始化SFTP客户端失败的日志
		Logger_caller("Failed to init sftp client: ", err,1)
		return nil,err
	}
	defer sftpClient.Close()

    // 打开指定路径的文件
	src_file,err := sftpClient.OpenFile(src,os.O_RDONLY)
	if err != nil{
		// 记录打开文件失败的日志
		Logger_caller("open file failed",err,1)
		return nil,err
	}
	defer src_file.Close()

    // 读取文件全部内容
	content,err := io.ReadAll(src_file)
	if err != nil{
		// 记录读取文件失败的日志
		Logger_caller("read file failed",err,1)
		return nil,err
	}
	
	return content,nil
}