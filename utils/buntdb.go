package utils

import (
	"fmt"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

// SetValue 是一个用于在buntDB数据库中设置键值对的函数
// 它接受一个*buntdb.DB类型的指针buntClient, 一个键key, 一个值value, 以及一个*zap.Logger类型的指针logger作为参数
// 如果设置键值对过程中发生错误, 该函数会记录错误信息并返回错误
func SetValue(buntClient *buntdb.DB, key, value string, logger *zap.Logger) error {
	// 使用buntClient的Update方法来执行数据库的写操作
	// Update方法接受一个函数作为参数, 该函数接收一个*buntdb.Tx类型的指针tx, 并返回一个错误（如果有）
	return buntClient.Update(func(tx *buntdb.Tx) error {
		// 尝试在事务tx中设置键key对应的值为value
		// 如果设置过程中没有发生错误, 则返回nil
		if _, _, err := tx.Set(key, value, nil); err != nil {
			// 如果设置键值对时发生错误, 使用logger记录错误信息
			// 错误信息包括尝试设置的键key, 值value, 以及错误详情
			logger.Error(fmt.Sprintf(`设置键值对"%s:%s"错误: [%s]`, key, value, err.Error()))
			// 返回错误, 终止操作
			return err
		}
		// 如果没有发生错误, 返回nil, 表示操作成功
		return nil
	})
}

// GetValue 从 buntClient 数据库中获取指定键的值
// 如果键不存在或者获取过程中出现错误, 将记录错误信息并返回错误
// 参数:
//
//	buntClient: buntdb 数据库的客户端连接
//	key: 需要获取值的键
//	logger: 用于记录错误信息的日志对象
//
// 返回值:
//
//	string: 键对应的值
//	error: 如果获取过程中出现错误, 则返回该错误
func GetValue(buntClient *buntdb.DB, key string, logger *zap.Logger) (string, error) {
	var value string

	// 使用 buntClient 的 View 方法以只读事务的形式获取键的值
	// View 方法的参数是一个函数, 该函数接收一个只读事务 tx, 并可能返回一个错误
	if err := buntClient.View(func(tx *buntdb.Tx) error {
		var err error
		// 从事务 tx 中获取键为 key 的值
		value, err = tx.Get(key)
		// 如果获取过程中出现错误, 记录错误信息并返回错误
		if err != nil {
			logger.Error(fmt.Sprintf(`获取键"%s"的值错误: [%s]`, key, err.Error()))
			return err
		}
		return nil
	}); err != nil {
		// 如果 View 方法执行失败, 返回空字符串和错误
		return "", err
	}

	// 返回获取到的值
	return value, nil
}

func GetKey(buntClient *buntdb.DB, value string, logger *zap.Logger) (string, error) {
	var key string
	const exist = false
	err := buntClient.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(index, val string) bool {
			if val == value {
				key = index
				return exist
			}
			return !exist
		})
		return nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf(`查找值"%s"的键出错: [%s]`, value, err.Error()))
		return "", err
	}
	return key, nil
}

// DeleteValue 从buntClient数据库中删除指定的键
// 如果键不存在, 函数将成功返回, 不会返回错误
// 参数:
//
//	buntClient: buntdb数据库的连接指针
//	key: 需要删除的键的字符串表示
//	logger: 用于记录错误信息的日志对象指针
//
// 返回值:
//
//	如果删除操作成功或键不存在, 则返回nil；否则返回错误
func DeleteValue(buntClient *buntdb.DB, key string, logger *zap.Logger) error {
	// 使用Update函数执行写操作, 它会确保数据库的一致性和并发安全
	err := buntClient.Update(func(tx *buntdb.Tx) error {
		// 尝试从事务中删除指定的键
		// 如果键不存在, Delete会返回一个错误
		if _, err := tx.Delete(key); err != nil {
			// 如果删除失败, 记录错误日志并返回错误
			logger.Error(fmt.Sprintf(`删除键"%s"的值出错: [%s]`, key, err.Error()))
			return err
		}
		// 如果删除成功, 返回nil表示操作成功
		return nil
	})
	// 如果Update函数返回错误, 将其返回给调用者
	if err != nil {
		return err
	}
	// 如果一切顺利, 返回nil表示操作成功
	return nil
}
