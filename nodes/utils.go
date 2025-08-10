package nodes

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sifu-box/ent"
	"sifu-box/ent/provider"
)

func updateNodes(name, uuid string, content []map[string]any, ent_client *ent.Client) error {
	data, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf(`"%s"出错, 序列化出站节点失败: [%s]`, name, err.Error())
	}
	if uuid != fmt.Sprintf("%x", md5.Sum([]byte(data))) {
		if err := ent_client.Provider.Update().Where(provider.NameEQ(name)).SetNodes(content).SetUUID(fmt.Sprintf("%x", md5.Sum([]byte(data)))).SetUpdated(true).Exec(context.Background()); err != nil {
			return fmt.Errorf(`"%s"出错, 保存数据失败: [%s]`, name, err.Error())
		}
	} else {
		if err := ent_client.Provider.Update().Where(provider.NameEQ(name)).SetUpdated(false).Exec(context.Background()); err != nil {
			return fmt.Errorf(`"%s"出错, 保存数据失败: [%s]`, name, err.Error())
		}
	}
	return nil
}
