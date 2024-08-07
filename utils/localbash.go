package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

// CommandExec 执行外部命令并返回其标准输出和标准错误输出
// 参数 command 是要执行的命令,args 是传递给命令的参数
// 返回值是命令的标准输出切片、标准错误输出切片和执行过程中可能发生的错误
func CommandExec(command string, args ...string) ([]string, []string, error) {
    // 使用 exec 包的 Command 函数启动外部命令
    cmd := exec.Command(command, args...)

    // 获取命令的标准输出管道
    stdoutPipe, err := cmd.StdoutPipe()
    if err != nil {
        // 如果获取标准输出管道失败,则返回错误
        return nil, nil, err
    }
    // 确保在函数返回前关闭标准输出管道
    defer stdoutPipe.Close()

    // 获取命令的标准错误输出管道
    errorsPipe, err := cmd.StderrPipe()
    if err != nil {
        // 如果获取标准错误输出管道失败,则返回错误
        return nil, nil, err
    }
    // 确保在函数返回前关闭标准错误输出管道
    defer errorsPipe.Close()

    // 启动命令
    if err := cmd.Start(); err != nil {
        // 如果启动命令失败,则返回错误
        return nil, nil, err
    }

    // 创建通道,用于在 goroutine 之间传递命令的输出和错误信息
    resultsCh := make(chan string)
    errorsCh := make(chan string)
    procErrs := make(chan error)
    var pipe sync.WaitGroup
    pipe.Add(2)

    // 启动一个 goroutine 来读取标准输出管道
    go func() {
        defer func() {
            pipe.Done()
            close(resultsCh)
        }()
        scanner := bufio.NewScanner(stdoutPipe)
        for scanner.Scan() {
            line := string(scanner.Bytes())
            resultsCh <- line
        }
        if scanner.Err() != nil {
            procErrs <- scanner.Err()
        }
    }()

    // 启动一个 goroutine 来读取标准错误输出管道
    go func() {
        defer func() {
            pipe.Done()
            close(errorsCh)
        }()
        scanner := bufio.NewScanner(errorsPipe)
        for scanner.Scan() {
            line := string(scanner.Bytes())
            errorsCh <- line
        }
        if scanner.Err() != nil {
            procErrs <- scanner.Err()
        }
    }()

    // 从结果和错误通道收集数据
    var results, errors []string
    for result := range resultsCh {
        results = append(results, result)
    }
    for msg := range errorsCh {
        errors = append(errors, msg)
    }

    // 等待两个 goroutine 完成
    pipe.Wait()
    close(procErrs)
    procErrsTag := false
    for proc_err := range procErrs {
        // 如果有处理错误发生,记录错误并设置标志
        LoggerCaller("没有EOF结尾标志", proc_err, 1)
        procErrsTag = true
    }

    // 根据处理错误标志和命令等待状态决定是否返回错误
    if procErrsTag {
        return results, errors, fmt.Errorf("get pipe output failed")
    }
    if err = cmd.Wait(); err != nil {
        return results, errors, err
    }
    return results, errors, nil
}