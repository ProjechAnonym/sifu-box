package utils

import (
	"fmt"
	"sync"

	"github.com/huandu/go-clone"
)

var globalVars = make(map[string]interface{})
var rlock sync.RWMutex

// GetValue 根据提供的键获取全局变量的值
// 它支持嵌套的键,这意味着可以通过一系列键获取嵌套的值
// 如果提供的键不存在或在解析过程中遇到问题,将返回错误
// 参数:
//   keys - 一个字符串键的可变长度列表,用于定位全局变量中的特定值
// 返回值:
//   interface{} - 请求的全局变量的克隆值,如果成功获取的话
//   error - 如果发生错误（如键不存在）,则返回错误信息
func GetValue(keys ...string) (interface{}, error) {
	
	// 读锁锁定以确保在读取全局变量时的一致性
	rlock.RLock()
	defer rlock.RUnlock()

	// 从全局变量开始,逐步解析键
	tempGlobalVars := globalVars
	for i, key := range keys {
		// 检查当前键对应的值是否存在
		if tempGlobalVars[key] != nil {
			// 如果已经是最后一个键,返回对应的值的克隆
			if i == len(keys)-1 {
				return clone.Clone(tempGlobalVars[key]), nil
			}
			// 如果当前值是一个映射,将其作为新的全局变量集合,继续解析
			if subMap, ok := tempGlobalVars[key].(map[string]interface{}); ok {
				tempGlobalVars = subMap
			} else {
				// 如果当前值不是映射,则返回错误,表示键不正确
				return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
			}
		} else {
			// 如果键不存在,则返回错误
			return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
		}
	}
	// 如果所有键都解析完毕,但没有找到对应的值,返回错误
	return nil, fmt.Errorf("参数不足,缺少键值参数")
}

// SetValue 函数用于将给定的值设置到全局变量中的特定键路径
// 它支持将值设置到嵌套的键路径中,如果路径中的任何键不存在或不是预期的类型,则返回错误
// 参数：
//   value: 要设置的值,类型为interface{},可以是任何类型
//   keys: 一个或多个字符串类型的键,用于指定值设置的路径
// 返回值：
//   如果设置成功,则返回nil；如果路径中的键不存在或类型不匹配,则返回错误
func SetValue(value interface{}, keys ...string) error {
	
	// 使用读锁来确保在读取和写入全局变量时的线程安全
	rlock.Lock()
	defer rlock.Unlock()

	// 从全局变量开始,将要设置的值沿着键的路径深入
	tempVars := globalVars
	
	// 遍历键的路径,检查并更新每个键对应的值
	for i, key := range keys {
		
		// 如果当前键是路径的最后一个键,则直接设置其值并返回成功
		if i == len(keys)-1 {
			tempVars[key] = value
			return nil
		} else {
			
			// 检查当前键的值是否为预期的map类型
			if sub_map, ok := tempVars[key].(map[string]interface{}); ok {
				
				// 如果是map类型,则继续深入到下一级
				tempVars = sub_map
			} else {
				
				// 如果当前键的值不是map类型,则返回错误
				return fmt.Errorf("参数%d '%s' 不存在", i+1, key)
			}
		}
	}
	
	// 如果遍历完键的路径后没有找到合适的设置位置,则返回参数不足的错误
	return fmt.Errorf("参数不足,缺少键值参数")
}



// DelValue 删除嵌套字典中的值
// 参数 keys 是一个字符串切片,表示要删除的键值路径
// 返回 error,如果路径中的某个键不是字典类型或路径不完整
func DelValue(keys ...string) error {
    // 加读锁以确保并发访问时的一致性
    rlock.Lock()
    defer rlock.Unlock()
    
    // 从全局变量中获取当前的环境变量
    tempVars := globalVars
    
    // 遍历键值路径,逐层定位到要删除的键
    for i, key := range keys {
        // 如果已经是最后一个键,则直接删除并返回
        if i == len(keys)-1 {
            delete(tempVars, key)
            return nil
        }
        
        // 尝试将当前值转换为字典类型
        if subMap, ok := tempVars[key].(map[string]interface{}); ok {
            // 转换成功,继续深入下一层
            tempVars = subMap
        } else {
            // 转换失败,返回错误信息
            return fmt.Errorf("参数%d '%s' 不存在", i+1, key)
        }
    }
    
    // 如果遍历结束还没删除,说明参数不足
    return fmt.Errorf("参数不足,缺少键值参数")
}
func Show(){
	fmt.Println(globalVars)
}