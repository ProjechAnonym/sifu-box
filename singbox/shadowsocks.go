package singbox

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"sifu-box/models"
	"strconv"
	"strings"
)

func MarshalShadowsocks(proxyMap map[string]interface{}) (map[string]interface{}, error) {
	
	ss := models.ShadowSocks{
		Type:        "shadowsocks",
		Tag:         proxyMap["name"].(string),
		Server:      proxyMap["server"].(string),
		Server_port: proxyMap["port"].(int),
		Method:      proxyMap["cipher"].(string),
		Password:    proxyMap["password"].(string),
	}

	
	
	ssMap, err := Struct2map(ss, "ss")
	if err != nil {
		return nil, err
	}

	
	return ssMap, nil
}
func Base64Shadowsocks(link string) (map[string]interface{}, error) {
    
	info, err := url.QueryUnescape(strings.TrimPrefix(link, "ss://"))
	if err != nil {
		return nil, err
	}

    
	parts := strings.Split(info, "@")
    
	decodedInfo, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

    
	infoParts := strings.Split(string(decodedInfo), ":")
    
	serverInfo := strings.Split(parts[1], "#")

    
	serverUrl, err := url.Parse("ss://" + serverInfo[0])
	if err != nil {
        
		return nil, fmt.Errorf("failed to parse server URL: %v", err)
	}

    
	port, err := strconv.Atoi(serverUrl.Port())
	if err != nil {
		return nil, err
	}

    
	ss := models.ShadowSocks{
		Type: "shadowsocks",
		Tag: serverInfo[1],
		Server: serverUrl.Hostname(),
		Server_port: port,
		Method: infoParts[0],
		Password: infoParts[1],
	}

    
	ssMap, err := Struct2map(ss, "shadowsocks")
	if err != nil {
        
		return nil, err
	}

    
	return ssMap, nil
}