package utils

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)
func Sftp_delete(server Server,path string) error{
	config := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host,err := url.Parse(server.Url)
	if err != nil{
		Logger_caller("parse url failed",err,1)
		return err
	}
	client, err := ssh.Dial("tcp", host.Hostname() + ":22", config)
	if err != nil {
		Logger_caller("Failed to dial: ", err,1)
	}
	defer client.Close()
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		Logger_caller("Failed to init sftp client: ", err,1)
	}
	defer sftpClient.Close()
	if _,err := sftpClient.Stat(path);err != nil{
		if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile){
			Logger_caller("file not found",err,1)
			return nil
		}else{
			Logger_caller("stat file failed",err,1)
			return err
		}
	}
	if err := sftpClient.RemoveAll(path);err != nil{
		Logger_caller("remove file failed",err,1)
		return err
	}
	return nil
}
func Sftp_write(server Server,content []byte, dst string) error{
	config := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host,err := url.Parse(server.Url)
	if err != nil{
		Logger_caller("parse url failed",err,1)
		return err
	}
	client, err := ssh.Dial("tcp", host.Hostname() + ":22", config)
	if err != nil {
		Logger_caller("Failed to dial: ", err,1)
	}
	defer client.Close()
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		Logger_caller("Failed to init sftp client: ", err,1)
	}
	defer sftpClient.Close()
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
	file, err := sftpClient.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
    defer func() {
        // 确保文件在函数返回前关闭,避免资源泄露
        if err := file.Close(); err != nil {
            Logger_caller("File can not close!", err,1)
        }
    }()
    if err != nil {
        Logger_caller("Create file failed", err,1)
        return err
    }
	// 将内容写入文件
	_, err = file.Write(content)
	if err != nil {
		Logger_caller("Write config failed!", err,1)
		return err
	}

	// 操作成功,返回nil
	return nil
}

func Sftp_read(server Server,src string) ([]byte,error){
	config := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host,err := url.Parse(server.Url)
	if err != nil{
		Logger_caller("parse url failed",err,1)
		return nil,err
	}
	client, err := ssh.Dial("tcp", host.Hostname() + ":22", config)
	if err != nil {
		Logger_caller("Failed to dial: ", err,1)
	}
	defer client.Close()
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		Logger_caller("Failed to init sftp client: ", err,1)
	}
	defer sftpClient.Close()
	src_file,err := sftpClient.OpenFile(src,os.O_RDONLY)
	if err != nil{
		Logger_caller("open file failed",err,1)
		return nil,err
	}
	defer src_file.Close()
	content,err := io.ReadAll(src_file)
	if err != nil{
		Logger_caller("read file failed",err,1)
		return nil,err
	}
	
	return content,nil
}