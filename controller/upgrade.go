package controller

import (
	"fmt"
	"io"
	"mime/multipart"
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
//   service: 需要升级的服务名称。
//   lock: 用于同步访问的互斥锁，以防止并发升级。
// 返回值:
//   如果升级过程中发生错误，则返回错误。
func UpgradeApp(file multipart.File, path, addr,service string,lock *sync.Mutex) error {
    // 查询数据库以获取目标主机的信息。
    var host models.Host
    if err := utils.DiskDb.Table("hosts").Where("url = ?", addr).First(&host).Error; err != nil {
        utils.LoggerCaller("数据库查询失败", err, 1)
        return fmt.Errorf("数据库查询失败")
    }

    // 尝试获取锁，以确保在同一时间内只有一个升级操作在执行。
    for {
        if lock.TryLock() {
            break
        }
    }
    defer lock.Unlock()

    // 读取文件内容。
    content,err := io.ReadAll(file)
    if err != nil {
        utils.LoggerCaller("文件写入失败", err, 1)
        return fmt.Errorf("文件写入失败")
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