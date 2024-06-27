package utils

import (
	"bufio"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)




func Command_ssh(server Server,command string,args ...string) ([]string,[]string,error){
	config := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host,err := url.Parse(server.Url)
	if err != nil{
		Logger_caller("parse server url failde",err,1)
		return nil,nil,err
	}
	client, err := ssh.Dial("tcp", host.Hostname() + ":22", config)
	if err != nil {
		Logger_caller("SSH connection failed",err,1)
		return nil,nil,err
	}
	defer client.Close()
	session,err := client.NewSession()
	if err != nil {
		Logger_caller("SSH session failed",err,1)
		return nil,nil,err
	}
	defer session.Close()
	// 重定向标准输出和标准错误
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

	results_ch_ssh := make(chan string)
	errors_ch_ssh := make(chan string)
	proc_errs_ssh := make(chan error)
	if err := session.Run(command + " " + strings.Join(args," ")); err != nil {
		Logger_caller("Failed to run command: ", err,1)
		return nil,nil,err
	}
	var ssh_pipe sync.WaitGroup
	ssh_pipe.Add(2)
    // 并发读取命令的标准输出
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
	// 并发读取命令的错误输出
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
	// 收集命令的标准输出和错误输出
	var results,errors []string
	for result := range results_ch_ssh {
		results = append(results,result)
	}
	for msg := range errors_ch_ssh {
		errors = append(errors,msg)
	}
    // 等待读取操作完成
	ssh_pipe.Wait()

	close(proc_errs_ssh)
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

    // 命令执行成功,返回输出结果
	return results,errors,nil
	
}