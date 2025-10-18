package control

import (
	"fmt"
	"sifu-box/application"
	"sifu-box/initial"
	"sifu-box/model"
	"sifu-box/utils"
	"time"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func FetchYacd(bunt_client *buntdb.DB, logger *zap.Logger) (*model.Yacd, error) {
	content, err := utils.GetValue(bunt_client, initial.YACD, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("获取用户配置信息失败: [%s]", err.Error())
	}
	yacd := model.Yacd{}
	if err := yaml.Unmarshal([]byte(content), &yacd); err != nil {
		logger.Error(fmt.Sprintf("序列化用户配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化用户配置信息失败: [%s]", err.Error())
	}
	return &yacd, nil
}
func SetTemplate(name string, bunt_client *buntdb.DB, signal_chan *chan application.Signal, web_chan *chan bool, logger *zap.Logger) error {
	if err := utils.SetValue(bunt_client, initial.ACTIVE_TEMPLATE, name, logger); err != nil {
		return fmt.Errorf("设置模板失败: [%s]", err.Error())
	}
	*signal_chan <- application.Signal{Cron: false, Operation: application.RELOAD_SERVICE}
	select {
	case res := <-*web_chan:
		if res {
			logger.Info(fmt.Sprintf(`模板切换"%s"成功`, name))
			return nil
		}
		logger.Error(fmt.Sprintf(`载入"%s"配置失败, sing-box服务异常`, name))
		return fmt.Errorf(`载入"%s"配置失败, sing-box服务异常`, name)
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		return fmt.Errorf(`接收操作结果超时`)
	}
}
