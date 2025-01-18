package utils

import (
	"fmt"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SetValue(buntClient *buntdb.DB, key, value string, logger *zap.Logger) error {
	return buntClient.Update(func(tx *buntdb.Tx) error {
		if _, _, err := tx.Set(key, value, nil); err != nil {
			logger.Error(fmt.Sprintf(`设置键值对"%s:%s"错误: [%s]`, key, value, err.Error()))
			return err
		}
		return nil
	})
}

func GetValue(buntClient *buntdb.DB, key string, logger *zap.Logger) (string, error) {
	var value string
	if err := buntClient.View(func(tx *buntdb.Tx) error {
		var err error
		value, err = tx.Get(key)
		if err != nil {
			logger.Error(fmt.Sprintf(`获取键"%s"的值错误: [%s]`, key, err.Error()))
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}
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

func DeleteValue(buntClient *buntdb.DB, key string, logger *zap.Logger) error {
	err := buntClient.Update(func(tx *buntdb.Tx) error {
		if _, err := tx.Delete(key); err != nil {
			logger.Error(fmt.Sprintf(`删除键"%s"的值出错: [%s]`, key, err.Error()))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}