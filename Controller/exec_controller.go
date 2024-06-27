package controller

import (
	"fmt"
	execute "sifu-box/Execute"
	singbox "sifu-box/Singbox"
	utils "sifu-box/Utils"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

// Update_config 根据给定的地址和配置更新服务器配置
// addr: 服务器地址,用于查找数据库中的服务器信息
// config: 配置字符串,用于更新服务器的配置
// lock: 互斥锁,用于保护应用程序启停
// 返回值: 错误信息,如果操作成功则为nil
func Update_config(addr, config string,lock *sync.Mutex) error {
    // 从数据库中查询服务器信息
    var server utils.Server
    if err := utils.Db.Model(&server).Select("localhost", "url", "username", "password").Where("url = ?", addr).First(&server).Error; err != nil {
        // 日志记录查询服务器信息失败
        utils.Logger_caller("search url failed!", err, 1)
        return err
    }

    // 获取代理配置中的代理信息
    proxy_config, err := utils.Get_value("Proxy")
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
func Refresh_items(lock *sync.Mutex) []error {
    // 配置工作流
    if errs := singbox.Config_workflow([]int{}); errs != nil {
        // 记录配置工作流失败的日志并返回错误信息
        return errs
    }

    // 从数据库中获取服务器列表
    var servers []utils.Server
    if err := utils.Db.Find(&servers).Error; err != nil {
        // 记录获取服务器失败的日志并返回错误信息
        utils.Logger_caller("Get servers failed", err, 1)
        return []error{fmt.Errorf("get servers failed")}
    }

    // 获取代理配置中的URL信息
    proxy_config, err := utils.Get_value("Proxy")
    if err != nil {
        // 记录获取代理配置失败的日志并返回错误
        utils.Logger_caller("load Proxy config failed", err, 1)
        return []error{err}
    }

    // 检查代理配置中的URL列表是否为空
    if len(proxy_config.(utils.Box_config).Url) == 0 {
        // 如果为空,记录错误日志并返回错误信息
        err := fmt.Errorf("no url in Proxy config")
        utils.Logger_caller("load Proxy url failed", err, 1)
        return []error{err}
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
    if errs := execute.Group_update(servers, proxy_config.(utils.Box_config), lock);errs != nil{
        return errs
    }

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
    var server utils.Server
    // 使用Gorm查询数据库,根据URL地址查找服务器信息
    if err := utils.Db.Model(&server).Where("url = ?", addr).First(&server).Error; err != nil {
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
        if (!strings.Contains(err.Error(), "exit status") && server.Localhost) || (!strings.Contains(err.Error(), "exited with status") && !server.Localhost){
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
    var server utils.Server
    if err := utils.Db.Model(&server).Where("url = ?", addr).First(&server).Error; err != nil {
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

// Set_interval 根据给定的时间间隔更新Cron任务的执行时间
// span: 代表时间间隔的不同维度,如分钟、小时、日期等
// cron_task: Cron任务对象,用于添加或删除Cron任务
// id: 当前Cron任务的ID,用于删除现有任务
// lock: 互斥锁,用于确保并发安全
// 返回值: 错误对象,如果操作失败则返回非nil的错误对象
func Set_interval(span []int,cron_task *cron.Cron,id *cron.EntryID,lock *sync.Mutex) error{
    // 根据span的长度生成不同的Cron表达式,用于控制任务执行的间隔
    var new_time string
    switch len(span) {
    case 0:
        new_time = ""
    case 1:
        new_time = fmt.Sprintf("*/%d * * * *",span[0])
    case 2:
        new_time = fmt.Sprintf("%d %d * * *",span[0],span[1])
    case 3:
        new_time = fmt.Sprintf("%d %d * * %d",span[0],span[1],span[2])
    }

    // 删除现有的Cron任务
    cron_task.Remove(*id)

    var err error
    // 如果新的Cron表达式不为空,则添加新的Cron任务
    if new_time != "" {
        *id,err = cron_task.AddFunc(new_time, func() {
            // 执行工作流配置的函数,此处未展示具体实现
            singbox.Config_workflow([]int{})
            var servers []utils.Server
            // 从数据库获取服务器列表
            if err := utils.Db.Find(&servers).Error; err != nil {
                // 记录获取服务器列表失败的日志
                utils.Logger_caller("get server list failed!", err, 1)
                return
            }
            // 获取代理配置
            proxy_config, err := utils.Get_value("Proxy")
            // 如果获取配置出错,记录错误信息
            if err != nil {
                // 记录获取代理配置失败的日志
                utils.Logger_caller("get proxy config failed", err, 1)
                return
            }
            // 使用互斥锁更新服务器组的代理配置
            execute.Group_update(servers, proxy_config.(utils.Box_config), lock)
        })
        if err != nil{
            // 记录添加Cron任务失败的日志
            utils.Logger_caller("set interval failed", err, 1)
            return err
        }
    }
    // 如果没有错误,则返回nil
    return nil
}