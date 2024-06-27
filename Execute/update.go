package execute

import (
	"io/fs"
	"path/filepath"

	utils "sifu-box/Utils"
)

// Backup_file 对给定的源文件进行备份
// 如果 localhost 参数为真,则备份发生在本地主机上,否则函数将简单地返回 nil
// 参数:
// origin - 源文件的路径
// backup - 备份文件应保存的路径
// perm - 备份文件的权限模式
// localhost - 指示备份是否应在本地主机上执行的布尔值
// 返回值:
// 如果备份过程中发生错误,则返回错误；否则返回 nil
func Backup_file(origin, backup string,perm fs.FileMode,server utils.Server) error{
    // 当备份应在本地主机上执行时
    if server.Localhost{
        // 尝试复制源文件到备份位置
        // 如果复制失败,记录错误并返回
        if err := utils.File_copy(origin,backup,perm); err != nil {
            utils.Logger_caller("copy original config file failed!",err,1)
            return err
        }
        // 在成功备份后,尝试删除原始文件
        // 如果删除失败,记录错误并返回
        if err := utils.File_delete(origin);err != nil{
            utils.Logger_caller("delete orginal config failed!",err,1)
            return err
        }
    }else{
        content,err := utils.Sftp_read(server,origin)
        if err != nil {
            utils.Logger_caller("copy original config file failed!",err,1)
            return err
        }
        if err := utils.File_write(content,backup,[]fs.FileMode{0755,perm});err != nil{
            utils.Logger_caller("copy original config file failed!",err,1)
            return err
        }
        if err := utils.Sftp_delete(server,origin);err != nil{
            utils.Logger_caller("delete orginal config failed!",err,1)
            return err
        }
    }
    // 备份过程成功,返回 nil 表示没有错误
    return nil
}
// Update_file 更新文件如果操作在本地服务器上,先创建备份目录,备份原文件,然后复制新文件到原文件位置
// origin_file: 原始文件路径
// new_file: 新文件路径
// backupfile: 备份文件路径
// perm: 新文件的权限模式
// server: 数据库服务器配置
// 返回错误信息,如果操作成功,则返回nil
func Update_file(origin_file, new_file, backupfile string, perm fs.FileMode, server utils.Server) error{
    // 创建备份文件所在的目录,权限模式为0755
    if err := utils.Dir_Create(filepath.Dir(backupfile),fs.FileMode(0755));err != nil{
        utils.Logger_caller("create backup dir failed!",err,1)
        return err
    }
    
    // 如果服务器是本地主机
    if server.Localhost{
        // 备份原文件
        if err := Backup_file(origin_file,backupfile,0644,server); err != nil {
            utils.Logger_caller("backup original config file failed!",err,1)
            return err
        }
        // 复制新文件到原始文件位置,设置新文件的权限模式
        if err := utils.File_copy(new_file,origin_file,perm); err != nil {
            utils.Logger_caller("copy new config file failed!",err,1)
            return err
        }
    }else{
        // 备份原文件
        if err := Backup_file(origin_file,backupfile,0644,server); err != nil {
            utils.Logger_caller("backup original config file failed!",err,1)
            return err
        }
        content,err := utils.File_read(new_file)
        if err != nil {
            utils.Logger_caller("read new config file failed!",err,1)
            return err
        }
        if err := utils.Sftp_write(server,content,origin_file);err != nil{
            utils.Logger_caller("copy new config file failed!",err,1)
            return err
        }
    }
    // 操作成功
    return nil
}
func Recover_file(origin_file,backup_file string,perm fs.FileMode, server utils.Server) error{
    if server.Localhost{
        // 复制新文件到原始文件位置,设置新文件的权限模式
        if err := utils.File_copy(backup_file,origin_file,perm); err != nil {
            utils.Logger_caller("copy new config file failed!",err,1)
            return err
        }
    }else{
        content,err := utils.File_read(backup_file)
        if err != nil {
            utils.Logger_caller("read new config file failed!",err,1)
            return err
        }
        if err := utils.Sftp_write(server,content,origin_file); err != nil {
            utils.Logger_caller("copy new config file failed!",err,1)
            return err
        }
    }
    return nil
}