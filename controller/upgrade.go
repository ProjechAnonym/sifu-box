package controller

import (
	"fmt"
	"os"
	"sifu-box/execute"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"
)

// UpgradeApp 升级应用程序。
// 该函数负责将应用程序的新版本部署到指定的主机上，并确保服务在升级过程中停止，在升级完成后重新启动。
// 参数:
//   file: 包含新版本应用程序的文件。
//   path: 应用程序文件在目标主机上的路径。
//   addr: 目标主机的地址。
//   service: 需要升级的服务名称。。
// 返回值:
//   如果升级过程中发生错误，则返回错误。
func upgradeApp(content []byte, path, addr, service string) error {
    // 查询数据库以获取目标主机的信息。
    var host models.Host
    if err := utils.DiskDb.Table("hosts").Where("url = ?", addr).First(&host).Error; err != nil {
        utils.LoggerCaller("数据库查询失败", err, 1)
        return fmt.Errorf("数据库查询失败")
    }

    // 停止目标主机上的服务。
    if err := execute.StopService(service, host); err != nil {
        utils.LoggerCaller("停止服务失败", err, 1)
        return fmt.Errorf("停止服务失败")
    }

    // 根据主机是否为本地主机，执行不同的文件操作。
    if host.Localhost {
        // 对于本地主机，删除旧的文件版本，写入新版本，并修改文件权限。
        if err := utils.FileDelete(path); err != nil {
            utils.LoggerCaller("文件删除失败", err, 1)
            return fmt.Errorf("文件删除失败")
        }
        if err := utils.FileWrite(content, path); err != nil {
            utils.LoggerCaller("文件写入失败", err, 1)
            return fmt.Errorf("文件写入失败")
        }
        if err := os.Chmod(path, 0755); err != nil {
            utils.LoggerCaller("文件权限修改失败", err, 1)
            return fmt.Errorf("文件权限修改失败")
        }
    }else{
        // 对于远程主机，通过SFTP删除旧的文件版本，写入新版本，并通过SSH修改文件权限。
        if err := utils.SftpDelete(host, path); err != nil {
            utils.LoggerCaller("文件删除失败", err, 1)
            return fmt.Errorf("文件删除失败")
        }
        if err := utils.SftpWrite(host, content, path); err != nil {
            utils.LoggerCaller("文件写入失败", err, 1)
            return fmt.Errorf("文件写入失败")
        }
        if _,_, err := utils.CommandSsh(host, "chmod", "0755", path);err != nil {
            utils.LoggerCaller("文件权限修改失败", err, 1)
            return fmt.Errorf("文件权限修改失败")
        }
    }

    // 检查服务状态，确保服务已经停止。
    status,err := execute.CheckService(service, host)
    if err != nil {
        utils.LoggerCaller("检查服务状态失败", err, 1)
        return fmt.Errorf("检查服务状态失败")
    }
    if !status {
        // 如果服务已经停止，则重新启动服务。
        if err := execute.BootService(service, host);err != nil{
            utils.LoggerCaller("启动服务失败", err, 1)
            return fmt.Errorf("启动服务失败")
        }
    }else{
        // 如果服务未正常停止，则返回错误。
        utils.LoggerCaller("服务未正常停止", fmt.Errorf("未知原因"), 1)
        return fmt.Errorf("服务未正常停止")
    }

    return nil
}

// UpgradeWorkflow 执行工作流升级。
// 该函数接收文件字节切片、地址切片、路径、服务名称和互斥锁作为参数，返回错误切片。
// 它并行地在给定地址上执行升级操作，并收集升级过程中的错误。
func UpgradeWorkflow(file []byte, addresses []string, path, service string, lock *sync.Mutex) []error {
    // 尝试获取锁，以确保同一时间只有一个升级流程在执行。
    for {
        if lock.TryLock() {
            break
        }
    }
    defer lock.Unlock()
    // upgradeTask 用于等待所有的更新操作完成
    var upgradeTask sync.WaitGroup
    // 初始化WaitGroup,添加主机数量加一的计数,用于等待所有更新完成
    upgradeTask.Add(len(addresses) + 1)

    // upgradeErrorsChan 用于收集更新过程中产生的错误
    upgradeErrorsChan := make(chan error, len(addresses))
    // upgradeErrors 用于存储所有的错误
    var upgradeErrors []error

    // upgradeCountChan 用于统计更新完成的主机数量
    upgradeCountChan := make(chan int, len(addresses))

    // 遍历主机列表,为每个主机启动一个更新操作
    for _, addr := range addresses {
        // 使用匿名协程执行更新操作
        go func() {
            // 更新完成后,减少WaitGroup的计数,并向countChan发送一个计数
            defer func() {
                upgradeTask.Done()
                upgradeCountChan <- 1
            }()

            // 执行更新操作,如果发生错误,则记录错误并发送到upgradeErrorsChan
            if err := upgradeApp(file, path, addr, service); err != nil {
                upgradeErrorsChan <- fmt.Errorf("%s升级失败", addr)
            }
        }()
    }
    
    // 使用匿名协程监控更新进度
    go func() {
        defer func() {
            // 更新完成后,减少WaitGroup的计数,并关闭upgradeErrorsChan和upgradeCountChan
            upgradeTask.Done()
            close(upgradeErrorsChan)
            close(upgradeCountChan)
        }()

        // sum 用于累计已完成的更新操作数量
        sum := 0
        // 遍历upgradeCountChan,统计更新完成的主机数量
        for count := range upgradeCountChan {
            sum += count
            // 当所有主机都已完成更新时,退出协程
            if sum == len(addresses) {
                return
            }
        }
    }()

    // 遍历upgradeErrorsChan,收集所有的错误
    for err := range upgradeErrorsChan {
        upgradeErrors = append(upgradeErrors, err)
    }

    // 等待所有的更新操作完成
    upgradeTask.Wait()
    // 返回所有的错误列表
    return upgradeErrors
}