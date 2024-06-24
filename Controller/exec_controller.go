package controller

import (
	database "sifu-box/Database"
	execute "sifu-box/Execute"
	utils "sifu-box/Utils"
	"sync"
)

// Update_config 根据给定的地址和配置更新服务器配置
// addr: 服务器地址,用于查找数据库中的服务器信息
// config: 配置字符串,用于更新服务器的配置
// lock: 互斥锁,用于保护应用程序启停
// 返回值: 错误信息,如果操作成功则为nil
func Update_config(addr, config string,lock *sync.Mutex) error {
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
    if err := execute.Exec_update(config, proxy_config.(utils.Box_config), server,false,lock);err != nil{
        // 日志记录更新配置失败
        utils.Logger_caller("update singbox config failed!", err, 1)
        return err
    }
    return nil
}