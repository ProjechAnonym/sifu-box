package controller

import (
	"fmt"
	"path/filepath"
	"sifu-box/execute"
	"sifu-box/models"
	"sifu-box/singbox"
	"sifu-box/utils"
	"sync"

	"gopkg.in/yaml.v3"
)

// FetchItems 从内存数据库中获取 Proxy 配置数据
// 该函数不接受任何参数
// 返回值是一个指向 models.Proxy 的指针和一个错误类型
// 如果在查询过程中发生错误,会记录错误日志并返回错误信息
func FetchItems() (*models.Proxy, error) {
    // 初始化 providers 切片,用于存储查询到的 Provider 数据
    var providers []models.Provider
    // 初始化 rulesets 切片,用于存储查询到的 Ruleset 数据
    var rulesets []models.Ruleset

    // 从内存数据库中查询 providers 数据
    // 如果查询失败,记录错误日志并返回错误信息
    if err := utils.MemoryDb.Find(&providers).Error; err != nil {
        utils.LoggerCaller("获取机场配置失败", err, 1)
        return nil, fmt.Errorf("获取机场配置失败")
    }

    // 从内存数据库中查询 rulesets 数据
    // 如果查询失败,记录错误日志并返回错误信息
    if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
        utils.LoggerCaller("获取规则集配置失败", err, 1)
        return nil, fmt.Errorf("获取规则集配置失败")
    }

    // 如果查询成功,返回包含查询结果的 Proxy 结构体指针和 nil 错误
    return &models.Proxy{Providers: providers, Rulesets: rulesets}, nil
}

// AddItems 添加新的代理配置项,并更新相关设置
// 参数:
// - newProxy: 待添加的代理配置
// - lock: 用于同步的互斥锁
// 返回值:
// - []error: 在执行过程中遇到的错误列表,如果没有错误则返回nil
func AddItems(newProxy models.Proxy, lock *sync.Mutex) []error {
    // 获取项目目录
    projectDir, err := utils.GetValue("project-dir")
    if err != nil {
        // 记录错误日志并返回
        utils.LoggerCaller("获取工作目录失败", err, 1)
        return []error{fmt.Errorf("获取工作目录失败")}
    }
	
    // 检查代理配置和规则集配置是否都为空
    if len(newProxy.Providers) == 0 && len(newProxy.Rulesets) == 0{
        // 如果都为空,则返回错误
        return []error{fmt.Errorf("没有有效配置")}
    }

    // 初始化用于存储添加过程中可能出现的错误的切片
	var addMsg []error

    // 如果有新的代理配置需要添加
	if len(newProxy.Providers) != 0 {
        // 尝试添加到数据库
		if err := utils.MemoryDb.Create(&newProxy.Providers).Error; err != nil {
			// 记录错误日志并添加到错误切片
			utils.LoggerCaller("添加机场配置失败", err, 1)
			addMsg = append(addMsg, fmt.Errorf("添加机场配置失败"))
		}
	}

    // 如果有新的规则集配置需要添加
	if len(newProxy.Rulesets) != 0 {
        // 尝试添加到数据库
		if err := utils.MemoryDb.Create(&newProxy.Rulesets).Error; err != nil {
			// 记录错误日志并添加到错误切片
			utils.LoggerCaller("添加规则集配置失败", err, 1)
			addMsg = append(addMsg,fmt.Errorf("添加规则集配置失败"))
		}
	}

    // 从数据库中获取最新的代理配置和规则集配置
	var newProviders []models.Provider
	var newRulesets []models.Ruleset
	if err := utils.MemoryDb.Find(&newProviders).Error; err != nil {
		// 记录错误日志并返回
		utils.LoggerCaller("获取机场配置失败", err, 1)
		return []error{fmt.Errorf("获取机场配置失败")}
    }
	if err := utils.MemoryDb.Find(&newRulesets).Error; err != nil {
		// 记录错误日志并返回
		utils.LoggerCaller("获取规则集配置失败", err, 1)
		return []error{fmt.Errorf("获取规则集配置失败")}
    }
    
    // 将获取到的配置序列化为yaml格式
	proxyYaml, err := yaml.Marshal(models.Proxy{Providers: newProviders,Rulesets: newRulesets})
	if err != nil {  
        // 记录错误日志并返回
		utils.LoggerCaller("序列化yaml文件失败", err, 1)
		return []error{fmt.Errorf("序列化yaml文件失败")}
	}

    // 将序列化后的配置写入文件
	if err := utils.FileWrite(proxyYaml, filepath.Join(projectDir.(string), "config", "proxy.config.yaml")); err != nil { 
        // 记录错误日志并返回
		utils.LoggerCaller("写入proxy配置文件失败", err, 1)
		return []error{fmt.Errorf("写入proxy配置文件失败")}
	}

    // 如果添加过程中有错误,返回错误切片
	if len(addMsg) != 0{
		return addMsg
    }

    // 尝试获取锁,以便安全地更新主机配置
	for {
		if lock.TryLock() {
			break
		}
	}
	defer lock.Unlock()

    // 从数据库中获取主机列表
	var hosts []models.Host
	if err := utils.DiskDb.Find(&hosts).Error; err != nil {
		// 记录错误日志并返回
		utils.LoggerCaller("获取主机列表失败", err, 1)
		return []error{fmt.Errorf("获取主机列表失败")}
	}

    // 初始化错误切片,用于存储执行流程中可能出现的错误
	var errs []error

    // 根据是否有规则集配置来决定执行的更新流程
	if len(newProxy.Rulesets) == 0 {
        // 如果没有规则集配置,则只更新指定的代理配置
		var specific []int
		for _,provider := range(newProxy.Providers){
			specific = append(specific,int(provider.Id))
		}
		errs = singbox.Workflow(specific...)
	} else {
        // 如果有规则集配置,则更新所有配置,并进行分组更新
		errs = singbox.Workflow()
		if len(errs) != 0 {
			return errs
		}
		errs = execute.GroupUpdate(hosts,newProviders,lock,false)
	}

    // 返回执行过程中遇到的错误,如果没有则返回nil
	return errs
}

// DeleteProxy 删除代理配置,包括提供者和规则集
// 参数 proxy 是一个映射,包含了待删除的提供者和规则集的标识
// 参数 lock 是一个互斥锁,用于确保并发安全
// 返回值是一个错误切片,包含了删除过程中可能发生的错误
func DeleteProxy(proxy map[string][]int,lock *sync.Mutex) []error{
    
    // 获取项目目录
    projectDir, err := utils.GetValue("project-dir")
    if err != nil {
        
        utils.LoggerCaller("获取工作目录失败", err, 1)
        return []error{fmt.Errorf("获取工作目录失败")}
    }
    // 初始化用于存储删除过程中出现的错误的切片
    var deletemsg []error
    // 如果有提供者待删除
	if len(proxy["providers"]) != 0 {
		// 初始化临时提供者切片,用于存储待删除的提供者
		var tempProviders []models.Provider
		var deleteProviders []models.Provider
		// 从内存数据库中查找待删除的提供者
		if err := utils.MemoryDb.Find(&tempProviders,proxy["providers"]).Error; err != nil {
			utils.LoggerCaller("获取待删除机场配置失败", err, 1)
			return []error{fmt.Errorf("获取待删除机场配置失败")}
		}
		// 遍历提供者,进行相关文件的删除操作
		for _,tempProvider := range(tempProviders){
			// 加密提供者名称为md5
			md5Label,err := utils.EncryptionMd5(tempProvider.Name)
			if err != nil {
				utils.LoggerCaller("加密md5失败",err,1)
				return []error{fmt.Errorf("加密md5失败")}
			}
			// 获取模板配置
			templates,err := utils.GetValue("templates")
			if err != nil {
				utils.LoggerCaller("获取模板配置失败", err, 1)
				return []error{fmt.Errorf("获取模板配置失败")}
			}
			// 遍历模板,删除相关配置文件
			var deleteTag bool
			for key := range(templates.(map[string]models.Template)){
				if err := utils.FileDelete(filepath.Join(projectDir.(string), "static", key, md5Label + ".json")); err != nil {
					utils.LoggerCaller(fmt.Sprintf("删除'%s'目录下的'%s'配置文件失败",key,tempProvider.Name),err,1)
					deleteTag = false
				}else{
					deleteTag = true
				}
			}
			// 如果提供者不是远程的,删除其yaml文件
			if !tempProvider.Remote{
				if err := utils.FileDelete(tempProvider.Path); err != nil {
					utils.LoggerCaller("删除yaml文件失败",err,1)
					deleteTag = false
				}else {
					deleteTag = true
				}
			}
			if deleteTag {
				deleteProviders = append(deleteProviders,tempProvider)
			}
        }
		// 从数据库中删除提供者配置
		if err := utils.MemoryDb.Delete(&deleteProviders).Error; err != nil {
			utils.LoggerCaller("删除机场配置失败", err, 1)
			deletemsg = append(deletemsg, fmt.Errorf("删除机场配置失败"))
		}
    }

    // 如果有规则集待删除
	if len(proxy["rulesets"]) != 0 {
		// 从数据库中删除规则集配置
		if err := utils.MemoryDb.Delete(&models.Ruleset{},proxy["rulesets"]).Error; err != nil {
			utils.LoggerCaller("删除规则集配置失败", err, 1)
			deletemsg = append(deletemsg, fmt.Errorf("删除规则集配置失败"))
		}
	}
    // 获取最新的提供者和规则集配置
	var providers []models.Provider
	var rulesets []models.Ruleset
	if err := utils.MemoryDb.Find(&providers).Error; err != nil {
		utils.LoggerCaller("获取机场配置失败", err, 1)
        return []error{fmt.Errorf("获取机场配置失败")}
    }
	if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
		utils.LoggerCaller("获取规则集配置失败", err, 1)
        return []error{fmt.Errorf("获取规则集配置失败")}
    }
	
    // 生成新的代理配置文件
	proxyYaml, err := yaml.Marshal(models.Proxy{Providers: providers,Rulesets: rulesets})
	if err != nil {
		
		utils.LoggerCaller("", err, 1)
		return []error{fmt.Errorf("解析代理配置文件失败")}
	}

	
	if err := utils.FileWrite(proxyYaml, filepath.Join(projectDir.(string), "config", "proxy.config.yaml")); err != nil {
		
		utils.LoggerCaller("生成代理配置文件失败", err, 1)
		return []error{fmt.Errorf("生成代理配置文件失败")}
	}
    // 如果删除过程中有错误,返回错误切片
	if len(deletemsg) != 0{
		return deletemsg
    }
    // 更新主机配置
	var hosts []models.Host
	if err := utils.DiskDb.Find(&hosts).Error; err != nil {
		utils.LoggerCaller("查询主机列表失败", err, 1)
		return []error{fmt.Errorf("查询主机列表失败")}
	}
    if len(proxy["providers"]) != 0 {
        // 更新主机的配置信息
        for _,host := range(hosts){
            changeTag := true
            if len(providers) == 0{
                changeTag = true
            }else{
                for _,provider := range(providers){
                    if host.Config == provider.Name{
                        changeTag = false
                        break
                    }
                }
            }
            if changeTag{
                if err := utils.DiskDb.Model(&models.Host{}).Where("url = ?",host.Url).Update("config","").Error; err != nil{
					utils.LoggerCaller("更换主机配置失败",err,1)
                    return []error{fmt.Errorf("更换主机配置失败")}
                }
            }
        }
    }
    // 执行singbox工作流和分组更新
	var errs []error
    if len(proxy["rulesets"]) != 0 {
		// 使用互斥锁确保并发安全
		for {
			if lock.TryLock() {
				break
			}
		}
		defer lock.Unlock()
		errs = singbox.Workflow()
		if len(errs) != 0 {
			return errs
		}
		errs = execute.GroupUpdate(hosts,providers,lock,false)
		if len(errs) != 0 {
			return errs
		}
    }
    // 如果没有错误,返回nil
    return nil
}