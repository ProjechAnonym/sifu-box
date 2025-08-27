package initial

import (
	"context"
	"fmt"

	"sifu-box/ent"

	"entgo.io/ent/dialect"
)

func InitEntdb(dir string) *ent.Client {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v\n", r)
		}
	}()
	entClient, err := ent.Open(dialect.SQLite, fmt.Sprintf("file:%s/sifu-box.db?cache=shared&_fk=1", dir))
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: [%s]", err.Error()))
	}
	if err = entClient.Schema.Create(context.Background()); err != nil {
		panic(fmt.Sprintf("创建表失败, 检查工作目录: [%s]", err.Error()))
	}
	return entClient
}
