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
		panic(fmt.Sprintf("Connecting the Database has failed: [%s]", err.Error()))
	}
	if err = entClient.Schema.Create(context.Background()); err != nil {
		panic(fmt.Sprintf("Creating Tables has failed, check the working diretory: [%s]", err.Error()))
	}
	return entClient
}
