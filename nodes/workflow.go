package nodes

import (
	"net/http"

	"go.uber.org/zap"
)

func Merge(logger *zap.Logger) {
	client := http.DefaultClient
	_, err := fetchFromRemote("test", "http://clashshare.cczzuu.top/node/20250803-clash.yaml", client, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
