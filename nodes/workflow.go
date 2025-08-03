package nodes

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func Merge(logger *zap.Logger) {
	// client := http.DefaultClient
	// ou, err := fetchFromRemote("test", "https://raw.githubusercontent.com/Pawdroid/Free-servers/main/sub", client, logger)
	ou, err := fetchFromLocal("test", "/opt/sifubox/1.yaml", logger)
	a, _ := json.MarshalIndent(ou, "", "  ")
	fmt.Println(string(a))
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
