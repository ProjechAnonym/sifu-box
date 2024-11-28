package execute

import (
	"fmt"
	"net/url"
	"path/filepath"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"
)

// ExecUpdate 执行配置更新流程
// 参数:
//   label - 配置标签
//   providers - 提供者列表,用于验证标签是否存在
//   host - 主机配置信息
//   specific - 指定是否为特定更新操作
//   lock - 用于控制并发访问的互斥锁
// 返回值:
//   error - 更新过程中若发生错误,返回该错误
func ExecUpdate(label string, providers []models.Provider, host models.Host, specific bool, lock *sync.Mutex) error {
	
	// 如果是特定更新操作,则尝试加锁以避免并发问题
	if specific {
		for {
			if lock.TryLock() {
				break
			}
		}
		defer lock.Unlock()
	}
	
	// 获取项目目录
	projectDir, err := utils.GetValue("project-dir")
	if err != nil {
		utils.LoggerCaller("获取工作目录失败", err, 1)
		return err
	}

	// 检查标签是否存在
	labelExist := false
	for _, proxy := range providers {
		if proxy.Name == label {
			labelExist = true
			break
		}
	}
	if !labelExist {
		return fmt.Errorf("标签'%s'不存在目前配置中", label)
	}

	// 生成新的配置文件名
	newFile, err := utils.EncryptionMd5(label)
	if err != nil {
		utils.LoggerCaller("MD5加密失败", err, 1)
		return err
	}

	// 解析主机URL
	link, err := url.Parse(host.Url)
	if err != nil {
		utils.LoggerCaller("主机url解析失败", err, 1)
		return err
	}
	backupFile := link.Hostname()

	// 定义原配置文件路径、备份文件路径和新配置文件路径
	originalPath := "/opt/singbox/config.json"
	backupPath := filepath.Join(projectDir.(string), "backup", backupFile+".json")
	newPath := filepath.Join(projectDir.(string), "static", host.Template, newFile+".json")
	// 更新配置文件
	if err := UpdateFile(originalPath, newPath, backupPath, host); err != nil {
		return err
	}

	// 重新加载配置
	if result, err := ReloadConfig("sing-box", host); err != nil || !result {
		
		// 若重新加载失败,尝试恢复原配置文件
		if recoverErr := RecoverFile(originalPath, backupPath, host); recoverErr != nil {
			return recoverErr
		}

		// 尝试重新启动服务
		if startErr := BootService("sing-box", host); startErr != nil {
			return startErr
		}

		return fmt.Errorf("reload new config failed")
	}

	// 更新数据库中的配置标签
	if err := utils.DiskDb.Model(&host).Where("url = ?", host.Url).Update("config", label).Error; err != nil {
		utils.LoggerCaller("更新数据库失败", err, 1)
		return err
	}
	
	// 记录更新成功的日志
	utils.LoggerCaller(fmt.Sprintf("更新'%s'成功,当前配置为: %s", host.Url, host.Config), nil, 1)
	return nil
}

// GroupUpdate 执行一组主机的更新操作
// 参数 hosts 是待更新的主机列表,providers 是服务提供者列表,lock 是用于同步的互斥锁,only 是控制是否仅允许一个更新操作的开关
// 返回值是更新过程中可能发生的错误列表
func GroupUpdate(hosts []models.Host, providers []models.Provider, lock *sync.Mutex, only bool) []error {
    // 如果only为true,表示需要独占更新操作
    if only {
        // 尝试获取锁,确保只有一个更新操作在执行
        for {
            if lock.TryLock() {
                break
            }
        }
        // 释放锁,确保其他操作可以在更新完成后继续
        defer lock.Unlock()
    }
    
    // hostsWorkflow 用于等待所有的更新操作完成
    var hostsWorkflow sync.WaitGroup
    // 初始化WaitGroup,添加主机数量加一的计数,用于等待所有更新完成
    hostsWorkflow.Add(len(hosts) + 1)
    
    // errChan 用于收集更新过程中产生的错误
    errChan := make(chan error, 3)
    // errList 用于存储所有的错误
    var errList []error
    
    // countChan 用于统计更新完成的主机数量
    countChan := make(chan int, 3)
    
    // 遍历主机列表,为每个主机启动一个更新操作
    for _, host := range hosts {
        // 使用匿名协程执行更新操作
        go func() {
            // 更新完成后,减少WaitGroup的计数,并向countChan发送一个计数
            defer func() {
                hostsWorkflow.Done()
                countChan <- 1
            }()
            
            // 执行更新操作,如果发生错误,则记录错误并发送到errChan
            if err := ExecUpdate(host.Config, providers, host, false, lock); err != nil {
                utils.LoggerCaller("update servers config failed", err, 1)
                errChan <- fmt.Errorf("主机'%s'配置'%s'更新失败", host.Url, host.Config)
            }
        }()
    }
    
    // 使用匿名协程监控更新进度
    go func() {
        defer func() {
            // 更新完成后,减少WaitGroup的计数,并关闭countChan和errChan
            hostsWorkflow.Done()
            close(countChan)
            close(errChan)
        }()
        
        // sum 用于累计已完成的更新操作数量
        sum := 0
        // 遍历countChan,统计更新完成的主机数量
        for count := range countChan {
            sum += count
            // 当所有主机都已完成更新时,退出协程
            if sum == len(hosts) {
                return
            }
        }
    }()
    
    // 遍历errChan,收集所有的错误
    for err := range errChan {
        errList = append(errList, err)
    }
    
    // 等待所有的更新操作完成
    hostsWorkflow.Wait()
    // 返回所有的错误列表
    return errList
}