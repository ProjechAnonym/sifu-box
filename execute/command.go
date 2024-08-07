package execute

import (
	"fmt"
	"sifu-box/models"
	"sifu-box/utils"
	"strings"
)

// ReloadConfig 重新加载指定服务的配置文件
// 参数 service 表示服务名称,host 表示主机信息
// 返回值为操作是否成功和可能的错误信息
func ReloadConfig(service string, host models.Host) (bool, error) {
    // 初始化最终状态为成功,后续根据执行结果进行修改
    finalStatus := true
    // 初始化用于存储命令执行结果和错误信息的切片
    var results, errors []string

    // 检查服务状态
    currentStatus, err := CheckService(service, host)
    if err != nil {
        // 如果服务检查失败,记录错误并返回
        utils.LoggerCaller(fmt.Sprintf("%s未运行", service), err, 1)
        return false, err
    }

    // 根据服务状态和是否为本地主机执行不同操作
    if host.Localhost {
        // 本地主机且服务正在运行时,尝试重新加载配置
        if currentStatus {
            _, _, err = utils.CommandExec("systemctl", "reload", service)
            if err != nil {
                // 如果重新加载配置失败,记录错误并返回
                utils.LoggerCaller("重载配置失败", err, 1)
                return false, err
            }
        } else {
            // 服务未运行时,尝试启动服务
            err = BootService(service, host)
            if err != nil {
                // 如果启动服务失败,记录错误并返回
                utils.LoggerCaller("启动服务失败", err, 1)
                return false, err
            }
        }
        // 获取服务日志以检查操作结果
        results, errors, err = utils.CommandExec("journalctl", "-u", service, "-n", "1")
        if err != nil {
            // 如果获取日志失败,记录错误并返回
            utils.LoggerCaller("获取日志文件失败", err, 1)
            return false, err
        }
    } else {
        // 远程主机且服务正在运行时,尝试通过SSH重新加载配置
        if currentStatus {
            _, _, err = utils.CommandSsh(host, "systemctl", "reload", service)
            if err != nil {
                // 如果通过SSH重新加载配置失败,记录错误并返回
                utils.LoggerCaller("重载配置失败", err, 1)
                return false, err
            }
        } else {
            // 服务未运行时,尝试通过SSH启动服务
            err = BootService(service, host)
            if err != nil {
                // 如果通过SSH启动服务失败,记录错误并返回
                utils.LoggerCaller("启动服务失败", err, 1)
                return false, err
            }
        }
        // 通过SSH获取服务日志以检查操作结果
        results, errors, err = utils.CommandSsh(host, "journalctl", "-u", service, "-n", "1")
        if err != nil {
            // 如果获取日志失败,记录错误并返回
            utils.LoggerCaller("获取日志文件失败", err, 1)
            return false, err
        }
    }

    // 检查日志结果,如果有错误信息,则更新最终状态
    for _, result := range results {
        if strings.Contains(result, "ERROR") {
            // 如果日志中包含错误信息,记录错误并更新最终状态为失败
            utils.LoggerCaller("重载配置失败", fmt.Errorf(result), 1)
            finalStatus = false
            break
        }
    }

    // 如果执行过程中有错误输出,则记录错误并返回错误
    if len(errors) != 0 {
        utils.LoggerCaller("错误", fmt.Errorf("命令出现错误返回"), 1)
        return false, fmt.Errorf("命令出现错误返回")
    }

    // 如果最终状态为失败,则返回失败信息
    if !finalStatus {
        return false, fmt.Errorf("重载新配置失败")
    }

    // 返回成功
    return finalStatus, nil
}

// BootService 用于启动指定的服务
// 参数service为要启动的服务名称,host为服务所在的主机信息
// 返回值为错误信息,如果没有错误则返回nil
func BootService(service string, host models.Host) error {
    // 初始化状态变量和错误变量
    var status bool
    var err error

    // 根据host的Localhost属性决定是本地启动服务还是通过SSH远程启动服务
    if host.Localhost {
        // 如果是本地主机,使用systemctl命令启动服务
        _, _, err = utils.CommandExec("systemctl", "start", service)
    } else {
        // 如果是远程主机,通过SSH执行systemctl命令启动服务
        _, _, err = utils.CommandSsh(host, "systemctl", "start", service)
    }

    // 如果启动服务时发生错误,记录错误并返回
    if err != nil {
        utils.LoggerCaller("启动服务失败", err, 1)
        return err
    }

    // 检查服务是否已经运行
    status, err = CheckService(service, host)
    if err != nil {
        // 如果检查服务状态时发生错误,记录错误并返回
        utils.LoggerCaller(fmt.Sprintf("%s未运行", service), err, 1)
        return err
    }

    // 如果服务状态为关闭,则返回错误
    if !status {
        return fmt.Errorf("%s状态为关闭", service)
    }

    // 服务启动成功,返回nil
    return nil
}
// StopService 用于停止指定的服务
// 参数:
//   service - 需要停止的服务名称
//   host - 操作的主机信息,包含是否是本地主机的标识
// 返回值:
//   如果停止服务时发生错误,返回该错误；否则返回nil
func StopService(service string, host models.Host) error {
    var err error
    // 根据是否是本地主机选择执行方式
    if host.Localhost {
        // 本地停止服务
        _, _, err = utils.CommandExec("systemctl", "stop", service)
    } else {
        // 远程停止服务
        _, _, err = utils.CommandSsh(host, "systemctl", "stop", service)
    }
    // 错误处理
    if err != nil {
        // 记录错误日志并返回错误
        utils.LoggerCaller("启动服务失败", err, 1)
        return err
    }
    // 服务停止成功,返回nil
    return nil
}
// CheckService 检查指定服务在给定主机上是否处于活动状态
// 它首先判断主机是否为本地主机,然后尝试通过本地命令执行或SSH来查询服务状态
// 如果服务是活动的,它返回true；否则返回false如果在查询过程中出现错误,它将返回一个错误
func CheckService(service string, host models.Host) (bool, error) {
    // 初始化状态变量为false
    status := false
    // 初始化用于存储命令执行结果和错误信息的切片
    var results, errors []string
    // err变量用于存储可能出现的错误
    var err error

    // 根据主机是否为本地主机决定使用本地命令执行还是SSH执行命令
    if host.Localhost {
        results, errors, err = utils.CommandExec("systemctl", "status", service)
    } else {
        results, errors, err = utils.CommandSsh(host, "systemctl", "status", service)
    }
    // 检查错误信息,如果错误信息中不包含特定的退出状态说明,返回错误
    if err != nil {
        if (!strings.Contains(err.Error(), "exit status") && host.Localhost) || (!strings.Contains(err.Error(), "exited with status") && !host.Localhost) {
            return false, err
        }
    }


    // 如果有命令执行错误,记录错误并返回false
    if len(errors) != 0 {
        utils.LoggerCaller("错误", fmt.Errorf("命令出现错误返回"), 1)
        return false, fmt.Errorf("命令出现错误返回")
    }

    // 遍历命令执行结果,查找服务活动状态
    for _, result := range results {
        // 如果发现服务是活动的,设置状态为true并结束遍历
        if strings.Contains(result, "active (running)") {
            status = true
            break
        }
    }

    // 返回服务是否活动的状态,如果没有错误,错误参数为nil
    return status, nil
}

