package execute

import (
	"io/fs"
	"path/filepath"

	utils "sifu-box/Utils"
)

// Backup_file 对给定的原始文件进行备份
// 如果服务器是本地主机,则直接复制原始文件到备份路径,并删除原始文件
// 如果服务器是远程主机,则通过SFTP协议下载原始文件,保存为备份文件,并在远程主机上删除原始文件
// 参数:
// origin - 原始文件的路径
// backup - 备份文件的路径
// perm - 文件的权限模式
// server - 远程服务器的信息
// 返回值:
// 如果操作过程中出现错误,则返回错误；否则返回nil
func Backup_file(origin, backup string, perm fs.FileMode, server utils.Server) error {
    // 检查是否在本地主机操作
    // 当备份应在本地主机上执行时
    if server.Localhost {
        // 本地文件复制
        if err := utils.File_copy(origin, backup, perm); err != nil {
            // 记录复制原始配置文件失败的日志
            utils.Logger_caller("copy original config file failed!", err, 1)
            return err
        }
        // 本地文件删除
        if err := utils.File_delete(origin); err != nil {
            // 记录删除原始配置文件失败的日志
            utils.Logger_caller("delete orginal config failed!", err, 1)
            return err
        }
    } else {
        // 从远程主机下载文件内容
        content, err := utils.Sftp_read(server, origin)
        if err != nil {
            // 记录从远程主机读取原始配置文件失败的日志
            utils.Logger_caller("copy original config file failed!", err, 1)
            return err
        }
        // 将远程文件内容写入本地备份文件
        if err := utils.File_write(content, backup, []fs.FileMode{0755, perm}); err != nil {
            // 记录写入备份文件失败的日志
            utils.Logger_caller("copy original config file failed!", err, 1)
            return err
        }
        // 通过SFTP协议在远程主机上删除原始文件
        if err := utils.Sftp_delete(server, origin); err != nil {
            // 记录删除原始配置文件失败的日志
            utils.Logger_caller("delete orginal config failed!", err, 1)
            return err
        }
    }
    // 操作成功,返回nil
    return nil
}

// Update_file 更新文件内容根据操作环境（本地或远程服务器）,直接覆盖原文件或通过SFTP协议上传新文件
// origin_file: 原始文件路径
// new_file: 新文件路径
// backupfile: 备份文件路径
// perm: 文件权限
// server: 服务器配置信息
// 返回值: 操作失败时返回错误信息,成功则返回nil
func Update_file(origin_file, new_file, backupfile string, perm fs.FileMode, server utils.Server) error{
    
    // 创建备份文件所在目录,确保目录存在并设置权限为0755
    if err := utils.Dir_Create(filepath.Dir(backupfile),fs.FileMode(0755));err != nil{
        utils.Logger_caller("创建备份目录失败！",err,1)
        return err
    }
  
    // 若在本地服务器上操作,执行文件备份及替换
    if server.Localhost{
        
        // 备份原始配置文件
        if err := Backup_file(origin_file,backupfile,0644,server); err != nil {
            utils.Logger_caller("backup origin file failed",err,1)
            return err
        }
        
        // 将新配置文件复制到原始文件位置,并设置文件权限
        if err := utils.File_copy(new_file,origin_file,perm); err != nil {
            utils.Logger_caller("copy new config file failed",err,1)
            return err
        }
    }else{
        // 若在远程服务器上操作,进行文件备份及上传
        // 备份原始配置文件
        if err := Backup_file(origin_file,backupfile,0644,server); err != nil {
            utils.Logger_caller("backup origin file failed",err,1)
            return err
        }
        // 读取新配置文件内容,准备上传至远程服务器
        content,err := utils.File_read(new_file)
        if err != nil {
            utils.Logger_caller("read new config file failed",err,1)
            return err
        }
        // 使用SFTP协议,将新配置文件内容上传至远程服务器的原始文件位置
        if err := utils.Sftp_write(server,content,origin_file);err != nil{
            utils.Logger_caller("upload new config file fail",err,1)
            return err
        }
    }
    
    // 文件更新完毕,返回nil表示成功
    return nil
}
// Recover_file 从备份文件恢复原始文件
// 如果服务器是本地主机,则直接复制备份文件到原始文件位置
// 如果服务器是远程主机,则先从备份文件读取内容,然后通过SFTP协议将内容写入远程主机的原始文件位置
// 参数:
// origin_file: 原始文件的路径
// backup_file: 备份文件的路径
// perm: 文件的权限模式
// server: 远程服务器的信息
// 返回值:
// 如果操作成功,则返回nil；如果操作失败,则返回相应的错误
func Recover_file(origin_file,backup_file string,perm fs.FileMode, server utils.Server) error{
    // 检查是否为本地主机
    if server.Localhost{
        // 对于本地主机,直接使用文件复制函数恢复原始文件
        if err := utils.File_copy(backup_file,origin_file,perm); err != nil {
            // 记录复制失败的日志并返回错误
            utils.Logger_caller("copy new config file failed!",err,1)
            return err
        }
    }else{
        // 对于远程主机,先从备份文件读取内容
        content,err := utils.File_read(backup_file)
        if err != nil {
            // 记录读取失败的日志并返回错误
            utils.Logger_caller("read new config file failed!",err,1)
            return err
        }
        // 使用SFTP协议将内容写入远程主机的原始文件位置
        if err := utils.Sftp_write(server,content,origin_file); err != nil {
            // 记录写入失败的日志并返回错误
            utils.Logger_caller("copy new config file failed!",err,1)
            return err
        }
    }
    // 操作成功,返回nil
    return nil
}