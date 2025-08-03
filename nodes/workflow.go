package nodes

import (
	"net/http"
	"sifu-box/ent"

	"go.uber.org/zap"
)

// vless、trojan、vmess节点
// https://raw.githubusercontent.com/Pawdroid/Free-servers/main/sub
func Merge(providers []*ent.Provider, logger *zap.Logger) {
	for _, provider := range providers {
		if provider.Remote {
			client := http.DefaultClient
			fetchFromRemote(provider.Name, provider.Path, client, logger)
		}
		fetchFromLocal(provider.Name, provider.Path, logger)
	}

}
