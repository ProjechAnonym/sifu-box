package controller

import (
	"fmt"
	database "sifu-box/Database"
	execute "sifu-box/Execute"
	singbox "sifu-box/Singbox"
	utils "sifu-box/Utils"
	"strings"
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

// Refresh_items 更新服务器配置项
// 使用锁 *sync.Mutex 来保证并发安全
// 返回错误信息如果更新过程中出现错误
func Refresh_items(lock *sync.Mutex) error {
    // 配置工作流
    if err := singbox.Config_workflow([]int{}); err != nil {
        // 记录配置工作流失败的日志并返回错误信息
        utils.Logger_caller("Config workflow failed", err, 1)
        return fmt.Errorf("config workflow failed")
    }

    // 从数据库中获取服务器列表
    var servers []database.Server
    if err := database.Db.Find(&servers).Error; err != nil {
        // 记录获取服务器失败的日志并返回错误信息
        utils.Logger_caller("Get servers failed", err, 1)
        return fmt.Errorf("get servers failed")
    }

    // 加载代理配置
    if err := utils.Load_config("Proxy"); err != nil {
        // 记录加载代理配置失败的日志并返回错误
        utils.Logger_caller("load Proxy config failed", err, 1)
        return err
    }

    // 获取代理配置中的URL信息
    proxy_config, err := utils.Get_value("Proxy")
    if err != nil {
        // 记录获取代理配置失败的日志并返回错误
        utils.Logger_caller("load Proxy config failed", err, 1)
        return err
    }

    // 检查代理配置中的URL列表是否为空
    if len(proxy_config.(utils.Box_config).Url) == 0 {
        // 如果为空,记录错误日志并返回错误信息
        err := fmt.Errorf("no url in Proxy config")
        utils.Logger_caller("load Proxy url failed", err, 1)
        return err
    }

    // 标记是否需要更新服务器配置
    server_update := false

    // 遍历服务器列表,检查是否需要更新配置
    for _, server := range servers {
        // 遍历代理配置中的URL,查找匹配的服务器配置
        for _, link := range proxy_config.(utils.Box_config).Url {
            // 如果找到匹配的配置,则标记需要更新并跳出循环
            if server.Config == link.Label {
                server_update = true
                break
            }
        }
        // 如果没有找到匹配的配置,将服务器配置更新为第一个URL的标签
        if !server_update {
            server.Config = proxy_config.(utils.Box_config).Url[0].Label
        }
    }

    // 执行服务器配置更新
    execute.Group_update(servers, proxy_config.(utils.Box_config), lock)

    // 更新完成,返回nil表示无错误
    return nil
}

// Check_status 根据服务地址和特定服务名称检查服务的状态
// 参数:
//   addr: 服务的URL地址
//   service: 需要检查的具体服务名称
// 返回值:
//   bool: 服务是否处于可用状态
//   error: 如果检查过程中出现错误,则返回错误信息
func Check_status(addr, service string) (bool, error) {
    // 从数据库中查询服务信息
    var server database.Server
    // 使用Gorm查询数据库,根据URL地址查找服务器信息
    if err := database.Db.Model(&server).Where("url = ?", addr).First(&server).Error; err != nil {
        // 如果查询出错,记录日志并返回错误
        utils.Logger_caller("get data fail", err, 1)
        return false, err
    }

    // 调用execute包中的Check_service函数检查指定服务的状态
    status, err := execute.Check_service(service, server)
    if err != nil {
        // 如果检查服务时出错,记录日志
        utils.Logger_caller("check service fail", err, 1)
        // 如果错误不包含"exit status",则返回错误；否则忽略错误,返回服务状态为false
        if !strings.Contains(err.Error(), "exit status") {
            return false, err
        }
    }
    // 返回服务状态
    return status, nil
}

// Boot_service 尝试启动指定服务
// 它首先从数据库中检索与给定地址相关联的服务器配置,
// 然后使用获取的配置尝试启动服务此函数设计为并发安全,
// 通过使用互斥锁来确保同一时间只启动一个服务实例
//
// 参数:
//   addr - 服务器地址,用于从数据库中查找服务器配置
//   service - 需要启动的服务名称
//   lock - 用于确保并发安全的互斥锁
//
// 返回值:
//   如果启动服务成功,则返回nil；否则返回相应的错误
func Boot_service(addr, service string, lock *sync.Mutex) error {
    // 从数据库中查找与给定地址匹配的服务器配置
    var server database.Server
    if err := database.Db.Model(&server).Where("url = ?", addr).First(&server).Error; err != nil {
        utils.Logger_caller("get data fail", err, 1)
        return err
    }

    // 使用互斥锁确保并发安全,尝试获取锁直到成功
    for {
        if lock.TryLock() {
            break
        }
    }
    defer lock.Unlock() // 在函数返回前释放锁

    // 使用获取的服务器配置尝试启动指定的服务
    if err := execute.Boot_service(service, server); err != nil {
        utils.Logger_caller("boot service fail", err, 1)
        return err
    }

    return nil
}