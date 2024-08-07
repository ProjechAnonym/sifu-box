package utils

import (
	"bufio"
	"fmt"
	"sifu-box/models"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

// CommandSsh 执行给定主机上的SSH命令。
// host: 目标主机的信息。
// command: 要执行的命令。
// args: 命令的参数。
// 返回值:
// - []string: 标准输出结果数组。
// - []string: 标准错误输出结果数组。
// - error: 错误信息，如果执行过程中出现错误。
func CommandSsh(host models.Host, command string, args ...string) ([]string, []string, error) {
    // 初始化SSH客户端配置。
    config, addr, err := InitClient(host)
    if err != nil {
        return nil, nil, err
    }
    
    // 基于配置和地址创建SSH客户端。
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        return nil, nil, err
    }
    
    // 确保客户端会话结束后关闭。
    defer client.Close()
    
    // 创建新的SSH会话。
    session, err := client.NewSession()
    if err != nil {
        return nil, nil, err
    }
    
    // 确保会话结束后关闭。
    defer session.Close()
    
    // 获取会话的标准输出和错误输出管道。
    stdout, err := session.StdoutPipe()
    if err != nil {
        return nil, nil, err
    }
    stderr, err := session.StderrPipe()
    if err != nil {
        return nil, nil, err
    }
    
    // 初始化结果、错误和处理错误的通道。
    resultsChSsh := make(chan string)
    errorsChSsh := make(chan string)
    procErrsSsh := make(chan error)
    
    // 启动SSH命令执行。
    if err := session.Run(command + " " + strings.Join(args, " ")); err != nil {
        return nil, nil, err
    }
    
    // 使用WaitGroup等待输出和错误读取完成。
    var sshPipe sync.WaitGroup
    sshPipe.Add(2)
    
    // 读取标准输出。
    go func() {
        defer func() {
            sshPipe.Done()
            close(resultsChSsh)
        }()
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            line := string(scanner.Bytes())
            resultsChSsh <- line
        }
        if scanner.Err() != nil {
            procErrsSsh <- scanner.Err()
        }
    }()
    
    // 读取标准错误输出。
    go func() {
        defer func() {
            sshPipe.Done()
            close(errorsChSsh)
        }()
        scanner := bufio.NewScanner(stderr)
        for scanner.Scan() {
            line := string(scanner.Bytes())
            errorsChSsh <- line
        }
        if scanner.Err() != nil {
            procErrsSsh <- scanner.Err()
        }
    }()
    
    // 收集标准输出和错误输出。
    var results, errors []string
    for result := range resultsChSsh {
        results = append(results, result)
    }
    for msg := range errorsChSsh {
        errors = append(errors, msg)
    }
    
    // 等待输出和错误读取完成。
    sshPipe.Wait()
    
    // 关闭错误处理通道。
    close(procErrsSsh)
    
    // 处理可能的扫描错误。
    procErrsTag := false
    for procErr := range procErrsSsh {
        LoggerCaller("没有EOF结束标志", procErr, 1)
        procErrsTag = true
    }
    
    // 如果有错误，返回错误。
    if procErrsTag {
        return results, errors, fmt.Errorf("获取命令输出失败")
    }
    
    // 返回执行命令的结果。
    return results, errors, nil
}