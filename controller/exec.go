package controller

import (
	"fmt"
	"sifu-box/execute"
	"sifu-box/models"
	"sifu-box/singbox"
	"sifu-box/utils"
	"sync"

	"github.com/robfig/cron/v3"
)

// UpdateConfig 更新主机配置
// 参数:
// addr - 主机地址,用于识别要更新的主机
// config - 新的配置内容
// lock - 用于同步的互斥锁
// 返回值:
// error - 如果更新过程中发生错误,返回相应的错误信息；否则返回nil
func UpdateConfig(addr, config string, lock *sync.Mutex) error {
    // 初始化一个Host对象,用于存储从数据库查询到的主机信息
    var host models.Host
    // 使用DiskDb查询数据库,获取指定地址的主机信息
    // 查询条件是主机的URL等于传入的addr
    // 如果查询错误,记录日志并返回错误信息
    if err := utils.DiskDb.Model(&host).Select("localhost", "url", "username", "password").Where("url = ?", addr).First(&host).Error; err != nil {
        utils.LoggerCaller("获取主机失败", err, 1)
        return fmt.Errorf("获取主机失败")
    }
    // 初始化一个Providers切片,用于存储从数据库查询到的所有代理提供者信息
    var providers []models.Provider
    // 使用MemoryDb查询数据库,获取所有代理提供者的信息
    // 如果查询错误,记录日志并返回错误信息
    if err := utils.MemoryDb.Find(&providers).Error; err != nil {
        utils.LoggerCaller("获取代理信息失败", err, 1)
        return fmt.Errorf("获取代理信息失败")
    }
    // 调用ExecUpdate函数,执行更新singbox配置的操作
    // 参数包括新的配置内容、代理提供者信息、目标主机信息,以及一个指示是否需要加锁的标志
    // 如果更新操作失败,记录日志并返回错误信息
    if err := execute.ExecUpdate(config, providers, host, true, lock); err != nil {
        utils.LoggerCaller("更新singbox配置失败", err, 1)
        return fmt.Errorf("更新singbox配置失败")
    }
    // 更新操作成功,返回nil
    return nil
}

// RefreshItems 刷新项目列表该函数主要负责更新主机和提供商的信息,并处理相关错误
// 参数 lock 用于确保线程安全,避免并发访问数据库时出现的问题
// 返回值为错误列表,如果执行过程中出现错误,会将错误信息放入列表中返回
func RefreshItems(lock *sync.Mutex) []error {
    // 使用锁前尝试立即获得锁,如果无法获得,则会阻塞直到锁可用
    for {
        if lock.TryLock() {
            break
        }
    }
    // 确保在函数返回前释放锁
    defer lock.Unlock()

    // 调用 singbox 的工作流,如果出现错误则直接返回错误列表
    if errs := singbox.Workflow(); errs != nil {
        return errs
    }

    // 初始化主机列表
    var hosts []models.Host
    // 从磁盘数据库中查询主机信息,如果出错则记录日志并返回错误列表
    if err := utils.DiskDb.Find(&hosts).Error; err != nil {
        utils.LoggerCaller("获取服务器失败", err, 1)
        return []error{fmt.Errorf("获取服务器失败")}
    }

    // 初始化提供商列表
    var providers []models.Provider
    // 从内存数据库中查询提供商信息,如果出错则记录日志并返回错误列表
    if err := utils.MemoryDb.Find(&providers).Error; err != nil {
        utils.LoggerCaller("获取代理信息失败", err, 1)
        return []error{fmt.Errorf("获取代理信息失败")}
    }

    // 如果提供商列表为空,则记录日志并返回错误
    if len(providers) == 0 {
        err := fmt.Errorf("配置中没有机场信息")
        utils.LoggerCaller("配置中没有机场信息", err, 1)
        return []error{err}
    }

    // 标记是否需要更新服务器配置
    serverUpdate := false

    // 遍历主机列表
    for _, host := range hosts {
        // 遍历提供商列表,检查主机配置是否需要更新
        for _, provider := range providers {
            if host.Config == provider.Name {
                serverUpdate = true
                break
            }
        }

        // 如果不需要更新,则为该主机分配第一个提供商的配置并更新数据库
        if !serverUpdate {
            host.Config = providers[0].Name
            if err := utils.DiskDb.Model(&models.Host{}).Where("url = ?", host.Url).Update("config", providers[0].Name).Error; err != nil {
                return []error{err}
            }
        }
    }

    // 如果主机列表不为空,则执行分组更新操作,如果出现错误则返回错误列表
    if len(hosts) != 0 {
        if errs := execute.GroupUpdate(hosts, providers, lock, false); errs != nil {
            return errs
        }
    }

    // 执行成功,返回空错误列表
    return nil
}
// CheckStatus 检查指定服务的状态
// 该函数首先根据地址查找主机信息,然后检查该主机上指定服务的状态
// 参数:
//   addr - 主机的地址,用于识别特定的主机
//   service - 要检查的服务名称
// 返回值:
//   bool - 服务是否可用的布尔值
//   error - 如果在查找主机或检查服务进程时发生错误,则返回错误
func CheckStatus(addr, service string) (bool, error) {
    // 初始化一个Host实例
    var host models.Host
    // 使用url作为条件查询数据库,以获取主机信息
    if err := utils.DiskDb.Model(&host).Where("url = ?", addr).First(&host).Error; err != nil {
        // 记录数据库查询错误
        utils.LoggerCaller("获取主机失败", err, 1)
        // 返回错误信息,指示获取主机失败
        return false, fmt.Errorf("获取主机失败")
    }
    // 调用CheckService函数检查服务状态
    status, err := execute.CheckService(service, host)
    if err != nil {
        // 记录服务进程检查失败
        utils.LoggerCaller("检查服务进程失败", err, 1)
        // 返回错误信息,指示检查服务进程失败
        return false, fmt.Errorf("检查服务进程失败")
    }
    // 返回服务状态,无错误
    return status, nil
}
// BootService 根据给定的地址和业务名称,启动一个服务
// 参数 addr 是主机的地址,service 是要启动的服务名称,lock 是用于同步的互斥锁
// 返回错误信息,如果启动服务失败
func BootService(addr, service string, lock *sync.Mutex) error {
    // 初始化一个 Host 结构体变量
    var host models.Host
    // 从数据库中查询主机信息
    if err := utils.DiskDb.Model(&host).Where("url = ?", addr).First(&host).Error; err != nil {
        // 记录日志并返回错误
        utils.LoggerCaller("获取主机失败", err, 1)
        return fmt.Errorf("获取主机失败")
    }

    // 使用互斥锁确保线程安全
    for {
        // 尝试获取锁,如果成功则跳出循环
        if lock.TryLock() {
            break
        }
    }
    // 确保在函数返回前释放锁
    defer lock.Unlock()

    // 调用 execute 包中的 BootService 函数启动服务
    if err := execute.BootService(service, host); err != nil {
        // 记录日志并返回错误
        utils.LoggerCaller("启动服务失败", err, 1)
        return fmt.Errorf("启动服务失败")
    }

    // 服务启动成功,返回 nil 表示没有错误
    return nil
}

// StopService 通过给定的地址和锁来停止指定的服务
// 参数:
//   addr - 主机的地址,用于识别要操作的主机
//   service - 要停止的服务名称
//   lock - 用于确保并发访问安全的互斥锁
// 返回值:
//   如果发生错误（如获取主机信息失败或停止服务失败）,则返回错误
func StopService(addr, service string, lock *sync.Mutex) error {
    // 初始化一个Host实例,用于后续的数据库查询
    var host models.Host
    // 使用地址查询数据库,获取对应的主机信息
    // 如果查询错误,则记录错误日志并返回错误
    if err := utils.DiskDb.Model(&host).Where("url = ?", addr).First(&host).Error; err != nil {
        utils.LoggerCaller("获取主机失败", err, 1)
        return fmt.Errorf("获取主机失败")
    }

    // 使用锁来控制并发,确保一次只有一个线程能执行到下面的代码
    for {
        if lock.TryLock() {
            break
        }
    }
    // 确保在函数返回前释放锁
    defer lock.Unlock()

    // 使用获取到的主机信息和指定的服务名来停止服务
    // 如果停止服务过程中出现错误,则记录错误日志并返回错误
    if err := execute.StopService(service, host); err != nil {
        utils.LoggerCaller("启动服务失败", err, 1)
        return fmt.Errorf("启动服务失败")
    }

    // 如果一切操作都成功,返回nil表示没有错误发生
    return nil
}
// SetInterval 根据给定的时间间隔设置定时任务
// span 参数表示时间间隔,cronTask 是cron的任务管理器,id是任务的ID,lock是用于同步的互斥锁
func SetInterval(span []int, cronTask *cron.Cron, id *cron.EntryID, lock *sync.Mutex) error {
    // 根据span的长度决定新的定时任务的时间格式
    var newTime string
    switch len(span) {
    case 0:
        newTime = ""
    case 1:
        newTime = fmt.Sprintf("*/%d * * * *",span[0])
    case 2:
        newTime = fmt.Sprintf("%d %d * * *",span[0],span[1])
    case 3:
        newTime = fmt.Sprintf("%d %d * * %d",span[0],span[1],span[2])
    }

    // 移除现有的定时任务
    cronTask.Remove(*id)

    var err error
    // 如果新的时间格式不为空,则添加新的定时任务
    if newTime != "" {
        // 添加定时任务,执行singbox的工作流
        *id,err = cronTask.AddFunc(newTime, func() {
            // 使用锁确保并发安全
            for {
                if lock.TryLock() {
                    break
                }
            }
            defer lock.Unlock()

            singbox.Workflow()

            // 获取主机和代理列表
            var hosts []models.Host
            var providers []models.Provider
            if err := utils.DiskDb.Find(&hosts).Error; err != nil {
                utils.LoggerCaller("获取主机列表失败", err, 1)
                return
            }
            if err := utils.MemoryDb.Find(&providers).Error; err != nil {
                utils.LoggerCaller("获取代理信息失败", err, 1)
                return
            }

            // 执行分组更新
            execute.GroupUpdate(hosts, providers, lock, false)
        })
        // 处理添加定时任务时可能出现的错误
        if err != nil{
            utils.LoggerCaller("修改定时任务失败", err, 1)
            return fmt.Errorf("修改定时任务失败")
        }
    }

    // 返回成功
    return nil
}