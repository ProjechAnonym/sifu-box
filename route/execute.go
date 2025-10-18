package route

import (
	"sifu-box/ent"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingExecute(api *gin.RouterGroup, user *model.User, bunt_client *buntdb.DB, ent_client *ent.Client, work_dir string, logger *zap.Logger) {

}
