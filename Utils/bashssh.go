package utils

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

// Command_ssh 通过SSH连接到服务器并执行指定的命令
// 参数:
//   server: 服务器配置信息
//   command: 待执行的命令
//   args: 命令的参数,可变长参数
// 返回值:
//   stdout的结果字符串数组
//   stderr的结果字符串数组
//   执行过程中可能出现的错误
func Command_ssh(server Server,command string,args ...string) ([]string,[]string,error){
    // 初始化SSH客户端配置并获取服务器地址
	config,addr,err := Init_sshclient(server)
	if err != nil {
		// 记录初始化SSH客户端时的错误
		Logger_caller("init ssh config failed",err,1)
		return nil,nil,err
	}
	// 建立SSH连接
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		// 记录SSH连接失败的错误
		Logger_caller("SSH connection failed",err,1)
		return nil,nil,err
	}
    // 确保SSH连接在函数返回前关闭
	defer client.Close()
	// 创建一个新的SSH会话
	session,err := client.NewSession()
	if err != nil {
		// 记录创建SSH会话失败的错误
		Logger_caller("SSH session failed",err,1)
		return nil,nil,err
	}
    // 确保SSH会话在函数返回前关闭
	defer session.Close()
    // 设置标准输出和标准错误的管道
	stdout, err := session.StdoutPipe()
	if err != nil {
		Logger_caller("Failed to stdout pipe: ", err,1)
		return nil,nil,err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		Logger_caller("Failed to stderr pipe: ", err,1)
		return nil,nil,err
	}

    // 创建通道以异步收集命令的标准输出和错误输出
	results_ch_ssh := make(chan string)
	errors_ch_ssh := make(chan string)
	proc_errs_ssh := make(chan error)
    // 执行命令
	if err := session.Run(command + " " + strings.Join(args," ")); err != nil {
		Logger_caller("Failed to run command: ", err,1)
		return nil,nil,err
	}
    // 同步读取标准输出和错误输出的等待组
	var ssh_pipe sync.WaitGroup
	ssh_pipe.Add(2)
    // 读取标准输出的协程
	go func ()  {
		defer func(){
			ssh_pipe.Done()
			close(results_ch_ssh)
		}()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			results_ch_ssh <- line
		}
		if scanner.Err() != nil {
			proc_errs_ssh <- scanner.Err()
		}
	}()
    // 读取标准错误的协程
	go func ()  {
		defer func(){
			ssh_pipe.Done()
			close(errors_ch_ssh)
		}()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			errors_ch_ssh <- line
		}
		if scanner.Err() != nil {
			proc_errs_ssh <- scanner.Err()
		}
	}()
    // 收集标准输出和错误输出的结果
	var results,errors []string
	for result := range results_ch_ssh {
		results = append(results,result)
	}
	for msg := range errors_ch_ssh {
		errors = append(errors,msg)
	}
    // 等待读取标准输出和错误输出的协程完成
	ssh_pipe.Wait()
    // 关闭错误通道
	close(proc_errs_ssh)
	// 检查读取过程中是否有错误发生
	proc_errs_tag := false
	for proc_err := range proc_errs_ssh {
        // 如果读取过程中有错误发生,记录错误
		Logger_caller("pipe without EOF tag",proc_err,1)
		proc_errs_tag = true
	}
    // 如果存在读取错误,返回错误信息
	if proc_errs_tag {
		return results,errors,fmt.Errorf("get pipe output failed")
	}

    // 返回命令执行的成功输出和错误输出
	return results,errors,nil
	
}