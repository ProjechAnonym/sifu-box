package initial

import (
	"fmt"
	"sifu-box/model"
	"sifu-box/utils"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func LoadSetting(config_path string, buntdb_client *buntdb.DB, logger *zap.Logger) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("%v\n", r)
	// 	}
	// }()
	content, err := utils.ReadFile(config_path)
	if err != nil {
		panic(fmt.Sprintf(`读取配置文件失败: [%s]`, err.Error()))
	}
	setting := struct {
		Smtp model.Smtp `json:"smtp" yaml:"smtp"`
		User model.User `json:"user" yaml:"user"`
		Yacd model.Yacd `json:"yacd" yaml:"yacd"`
	}{}
	if err := yaml.Unmarshal(content, &setting); err != nil {
		panic(fmt.Sprintf(`序列化配置文件失败: [%s]`, err.Error()))
	}
	smtp_content, err := yaml.Marshal(setting.Smtp)
	if err != nil {
		panic(fmt.Sprintf(`序列化SMTP配置失败: [%s]`, err.Error()))
	}
	if err := utils.SetValue(buntdb_client, SMTP, string(smtp_content), logger); err != nil {
		panic(fmt.Sprintf(`保存SMTP配置失败: [%s]`, err.Error()))
	}
	user_content, err := yaml.Marshal(setting.User)
	if err != nil {
		panic(fmt.Sprintf(`序列化SMTP配置失败: [%s]`, err.Error()))
	}
	if err := utils.SetValue(buntdb_client, USER, string(user_content), logger); err != nil {
		panic(fmt.Sprintf(`保存SMTP配置失败: [%s]`, err.Error()))
	}
	yacd_content, err := yaml.Marshal(setting.Yacd)
	if err != nil {
		panic(fmt.Sprintf(`序列化Yacd配置失败: [%s]`, err.Error()))
	}
	if err := utils.SetValue(buntdb_client, YACD, string(yacd_content), logger); err != nil {
		panic(fmt.Sprintf(`保存Yacd配置失败: [%s]`, err.Error()))
	}
}
