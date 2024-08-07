package execute

import (
	"path/filepath"
	"sifu-box/models"
	"sifu-box/utils"
)

// BackupFile 备份指定文件如果指定的主机是本地主机,则直接在本地进行操作；否则通过SFTP进行远程操作
// 参数:
//   origin: 原始文件路径
//   backup: 备份文件路径
//   host: 操作的主机信息
// 返回值:
//   error: 备份过程中出现的错误,如果没有错误则返回nil
func BackupFile(origin, backup string, host models.Host) error {
    // 判断是否为本地主机
    if host.Localhost {
        // 本地复制文件
        if err := utils.FileCopy(origin, backup); err != nil {
            utils.LoggerCaller("复制文件失败", err, 1)
            return err
        }
        // 删除原始文件
        if err := utils.FileDelete(origin); err != nil {
            utils.LoggerCaller("删除文件失败", err, 1)
            return err
        }
    } else {
        // 从远程主机读取文件内容
        content, err := utils.SftpRead(host, origin)
        if err != nil {
            utils.LoggerCaller("读取远程文件失败", err, 1)
            return err
        }
        // 在本地写入文件内容
        if err := utils.FileWrite(content, backup); err != nil {
            utils.LoggerCaller("写入本地文件失败", err, 1)
            return err
        }
        // 在远程主机删除原始文件
        if err := utils.SftpDelete(host, origin); err != nil {
            utils.LoggerCaller("删除远程文件失败", err, 1)
            return err
        }
    }
    return nil
}

// UpdateFile 更新文件如果主机是本地的,则直接替换文件；否则,将新文件内容上传到远程主机
// 参数:
// originFile - 原始文件路径
// newFile - 新文件路径
// backupFile - 备份文件路径
// host - 主机信息,包含是否是本地主机以及远程连接信息
// 返回值:
// error - 如果执行过程中出现错误,返回该错误
func UpdateFile(originFile, newFile, backupFile string, host models.Host) error{
    // 创建备份目录,确保备份文件的路径存在
    if err := utils.DirCreate(filepath.Dir(backupFile));err != nil{
        utils.LoggerCaller("创建备份目录失败！",err,1)
        return err
    }
    // 根据是否是本地主机,执行不同的文件更新策略
    if host.Localhost{
        // 备份原文件,以防更新失败需要恢复
        if err := BackupFile(originFile,backupFile,host); err != nil {
            utils.LoggerCaller("备份原文件失败",err,1)
        }
        // 使用新文件替换原文件
        if err := utils.FileCopy(newFile,originFile); err != nil {
            utils.LoggerCaller("设置新配置文件失败",err,1)
            return err
        }
    }else{
        // 对于远程主机,先备份原文件
        if err := BackupFile(originFile,backupFile,host); err != nil {
            utils.LoggerCaller("备份原文件失败",err,1)
        }
        // 读取新文件内容,准备上传到远程主机
        content,err := utils.FileRead(newFile)
        if err != nil {
            utils.LoggerCaller("设置新配置文件失败",err,1)
            return err
        }
        // 将新文件内容上传到远程主机
        if err := utils.SftpWrite(host,content,originFile);err != nil{
            utils.LoggerCaller("上传新配置文件到远程服务器失败",err,1)
            return err
        }
    }
    // 执行到这里说明更新成功,返回nil表示没有错误
    return nil
}

// RecoverFile 用于恢复原始文件当主机是本地主机时,它通过复制备份文件来恢复；否则,它通过读取备份文件内容并通过SFTP写入原始文件来恢复
// 参数:
// origin_file - 原始文件路径
// backup_file - 备份文件路径
// host - 主机信息,包含是否是本地主机以及SFTP连接信息等
// 返回值:
// error - 如果恢复过程中发生错误,返回相应的错误；否则返回nil
func RecoverFile(origin_file, backup_file string, host models.Host) error {
    // 当主机是本地主机时,直接通过复制文件的方式恢复
    if host.Localhost {
        if err := utils.FileCopy(backup_file, origin_file); err != nil {
            // 如果复制失败,记录错误并返回
            utils.LoggerCaller("恢复原配置文件失败", err, 1)
            return err
        }
    } else {
        // 当主机不是本地主机时,通过读取备份文件内容,然后通过SFTP写入远程主机的原始文件
        content, err := utils.FileRead(backup_file)
        if err != nil {
            // 如果读取备份文件内容失败,记录错误并返回
            utils.LoggerCaller("读取备份文件内容失败", err, 1)
            return err
        }
        
        // 使用SFTP将备份文件内容写入远程主机的原始文件
        if err := utils.SftpWrite(host, content, origin_file); err != nil {
            // 如果写入失败,记录错误并返回
            utils.LoggerCaller("写入远程主机原配置文件内容失败", err, 1)
            return err
        }
    }
    
    // 如果一切操作顺利,返回nil
    return nil
}