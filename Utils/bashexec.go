package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

// Command_exec 执行外部命令,并返回命令的标准输出和错误输出
// command: 要执行的命令名称
// args: 命令的参数,可变参数
// 返回值:
//   []string: 命令的标准输出
//   []string: 命令的错误输出
//   error: 执行过程中可能出现的错误
func Command_exec(command string,args ...string) ([]string,[]string,error){
    // 创建一个命令实例
	cmd := exec.Command(command,args...)
    
    // 创建一个标准输出的管道
	stdout_pipe, err := cmd.StdoutPipe()
	if err != nil {
        // 如果创建标准输出管道失败,记录错误并返回
		Logger_caller("failed to create stdout pipe",err,1)
		return nil,nil,err
	}
    // 确保在函数返回时关闭标准输出管道
	defer stdout_pipe.Close()

    // 创建一个错误输出的管道
	errors_pipe, err := cmd.StderrPipe()
	if err != nil {
        // 如果创建错误输出管道失败,记录错误并返回
		Logger_caller("failed to create error pipe",err,1)
		return nil,nil,err
	}
    // 确保在函数返回时关闭错误输出管道
	defer errors_pipe.Close()

    // 启动命令的执行
	if err := cmd.Start(); err != nil {
        // 如果启动命令失败,记录错误并返回
		Logger_caller("failed to create error pipe",err,1)
		return nil,nil,err
	}

    // 创建通道用于接收命令的标准输出和错误输出
	results_ch := make(chan string)
	errors_ch := make(chan string)
	proc_errs := make(chan error)
	var pipe sync.WaitGroup
	pipe.Add(2)

    // 并发读取命令的标准输出
	go func ()  {
		defer func(){
			pipe.Done()
			close(results_ch)
		}()
		scanner := bufio.NewScanner(stdout_pipe)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			results_ch <- line
		}
		if scanner.Err() != nil {
			proc_errs <- scanner.Err()
		}
	}()

    // 并发读取命令的错误输出
	go func ()  {
		defer func(){
			pipe.Done()
			close(errors_ch)
		}()
		scanner := bufio.NewScanner(errors_pipe)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			errors_ch <- line
		}
		if scanner.Err() != nil {
			proc_errs <- scanner.Err()
		}
	}()

    // 收集命令的标准输出和错误输出
	var results,errors []string
	for result := range results_ch {
		results = append(results,result)
	}
	for msg := range errors_ch {
		errors = append(errors,msg)
	}

    // 等待读取操作完成
	pipe.Wait()
	close(proc_errs)
	proc_errs_tag := false
	for proc_err := range proc_errs {
        // 如果读取过程中有错误发生,记录错误
		Logger_caller("pipe without EOF tag",proc_err,1)
		proc_errs_tag = true
	}
    // 如果存在读取错误,返回错误信息
	if proc_errs_tag {
		return results,errors,fmt.Errorf("get pipe output failed")
	}

    // 等待命令执行完成,并检查是否有错误发生
	if err = cmd.Wait();err != nil {
        // 如果命令执行失败,记录错误并返回
		Logger_caller("exec command failed!",err,1)
		return results,errors,err
	}
    // 命令执行成功,返回输出结果
	return results,errors,nil
}