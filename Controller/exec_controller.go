package controller

import (
	"fmt"
	"net/url"
	"path/filepath"
	database "sifu-box/Database"
	execute "sifu-box/Execute"
	utils "sifu-box/Utils"
)

// Update_config 根据给定的地址和配置更新服务器配置
// addr: 服务器地址,用于查找数据库中的服务器信息
// config: 配置字符串,用于更新服务器的配置
// 返回值: 错误信息,如果操作成功则为nil
func Update_config(addr, config string) error {
    // 获取项目目录,用于后续备份配置文件
    project_dir, err := utils.Get_value("project-dir")
    if err != nil {
        // 日志记录获取项目目录失败
        utils.Logger_caller("get project dir failed!", err, 1)
        return err
    }

    // 从数据库中查询服务器信息
    var server database.Server
    if err := database.Db.Model(&server).Select("localhost", "url", "username", "password").Where("url = ?", addr).First(&server).Error; err != nil {
        // 日志记录查询服务器信息失败
        utils.Logger_caller("search url failed!", err, 1)
        return err
    }

    // 加载代理配置
    if err := utils.Load_config("Proxy"); err != nil {
        // 日志记录加载代理配置失败
        utils.Logger_caller("load Proxy config failed!", err, 1)
        return err
    }

    // 获取代理配置中的代理信息
    proxy_config, err := utils.Get_value("Proxy")
    // 释放代理配置资源
    // 结束后删除代理信息配置
    defer utils.Del_key("Proxy")
    if err != nil {
        // 日志记录获取代理配置失败
        utils.Logger_caller("get Proxy config failed!", err, 1)
        return err
    }

    // 遍历代理配置,查找与config匹配的标签
    var label, backupfile string
    for _, proxy := range proxy_config.(utils.Box_config).Url {
        if proxy.Label == config {
            // 对配置字符串进行MD5加密
            label, err = utils.Encryption_md5(config)
            if err != nil {
                // 日志记录加密失败
                utils.Logger_caller("encryption md5 failed!", err, 1)
                return err
            }

            // 解析服务器URL,用于生成备份文件名
            link, err := url.Parse(server.Url)
            if err != nil {
                // 日志记录URL解析失败
                utils.Logger_caller("parse url failed", err, 1)
                return err
            }
            backupfile = link.Hostname()
            break
        }
    }

    // 检查是否找到匹配的标签,未找到则返回错误
    if label == "" {
        err := fmt.Errorf("%s label mismatch the urls in the proxy config", config)
        // 日志记录标签匹配失败
        utils.Logger_caller("label is empty!", err, 1)
        return err
    } else {
        // 构建原始配置文件、备份配置文件和新配置文件的路径
        origin_config := filepath.Join("/","opt","singbox","config.json")
        backup_config := filepath.Join(project_dir.(string),"temp","configbackup",fmt.Sprintf("%s.json",backupfile))
        new_config := filepath.Join(project_dir.(string),"static","Default",fmt.Sprintf("%s.json",label))

        // 更新配置文件
        if err := execute.Update_file(origin_config, new_config, backup_config, 0644, server); err != nil {
            // 日志记录更新配置文件失败
            utils.Logger_caller("update config failed!", err, 1)
            return err
        }   
    }
    // 查看更新配置后是否运行成功
    result,err := execute.Reload_config(server)
    if err != nil {
        // 日志记录重启配置文件失败
        utils.Logger_caller("reload config failed!", err, 1)
        // 恢复备份的配置文件
        if rec_err := execute.Recover_file(filepath.Join("/","opt","singbox","config.json"), filepath.Join(project_dir.(string),"temp","configbackup",fmt.Sprintf("%s.json",backupfile)), 0644, server);rec_err != nil{
            return rec_err
        }
        // 出现错误,尝试重启服务
        if start_err := execute.Boot_service("sing-box",server);start_err != nil{
            return start_err
        }
        return err
    }
    // 如果没有错误但是服务没有重载,尝试重启服务
    if !result {
        // 恢复备份的配置文件
        if rec_err := execute.Recover_file(filepath.Join("/","opt","singbox","config.json"), filepath.Join(project_dir.(string),"temp","configbackup",fmt.Sprintf("%s.json",backupfile)), 0644, server);rec_err != nil{
            return rec_err
        }
        if start_err := execute.Boot_service("sing-box",server);start_err != nil{
            return start_err
        }
        return fmt.Errorf("reload config failed")
    }

    // 更新数据库中服务器的配置信息
    if err := database.Db.Model(&server).Where("url = ?", addr).Update("config", config).Error; err != nil {
        // 日志记录更新数据库失败
        utils.Logger_caller("update server config failed!", err, 1)
        return err
    }

    // 操作成功,返回nil
    return nil
}