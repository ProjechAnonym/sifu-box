package nodes

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func Merge(logger *zap.Logger) {
	client := http.DefaultClient
	ou, err := fetchFromRemote("test", "https://sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d", client, logger)
	fmt.Println(ou)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
