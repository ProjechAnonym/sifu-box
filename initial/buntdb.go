package initial

import (
	"fmt"

	"github.com/tidwall/buntdb"
)

func InitBuntdb() *buntdb.DB {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v\n", r)
		}
	}()
	buntClient, err := buntdb.Open(":memory:")
	if err != nil {
		panic(fmt.Sprintf(`创建内存数据库buntdb失败: [%s]`, err.Error()))
	}
	return buntClient
}
