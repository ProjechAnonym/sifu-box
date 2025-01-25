package initial

import (
	"fmt"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func InitBuntdb(logger *zap.Logger) *buntdb.DB{
	buntClient, err := buntdb.Open(":memory:")
	if err != nil {
		logger.Error(fmt.Sprintf("连接Buntdb数据库失败: [%s]",err.Error()))
		panic(err)
	}
	return buntClient
}