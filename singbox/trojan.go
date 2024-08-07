package singbox

import (
	"net/url"
	"sifu-box/models"
	"strconv"
	"strings"
)

func MarshalTrojan(proxyMap map[string]interface{}) (map[string]interface{}, error) {
	
	skipCertVerify, err := GetMapValue(proxyMap, "skip-cert-verify")
	if err != nil {
		skipCertVerify = false
	}

	sni, err := GetMapValue(proxyMap, "sni")
	if err != nil {
		return nil, err
	}

	trojan := models.Trojan{
		Type:        "trojan",
		Tag:         proxyMap["name"].(string),
		Server:      proxyMap["server"].(string),
		Server_port: proxyMap["port"].(int),
		Password:    proxyMap["password"].(string),
		Tls: &models.Tls{
			Enabled:     true,
			Insecure:    skipCertVerify.(bool),
			Server_name: sni.(string),
		},
	}
	trojanMap, err := Struct2map(trojan, "trojan")
	if err != nil {
		return nil, err
	}
	return trojanMap, nil
}

func Base64Trojan(link string) (map[string]interface{}, error) {
    
    info := strings.TrimPrefix(link, "trojan://")
    
    parts := strings.Split(info, "@")
    
    password := parts[0]
    
    urlParts := strings.Split(parts[1], "#")
    
    serverUrl, err := url.Parse("trojan://" + urlParts[0])
    if err != nil {       
        return nil, err
    }
    
    tag, err := url.QueryUnescape(urlParts[1])
    if err != nil {
        return nil, err
    }
    
    port, err := strconv.Atoi(serverUrl.Port())
    if err != nil {
        return nil, err
    }
    
    
    params := serverUrl.Query()
    
    var skipCert bool
    
    if skipCertVerify := params.Get("allowInsecure"); skipCertVerify != "" {
        if skipCertVerify == "1" {
            skipCert = true
        } else {
            skipCert = false
        }
    } else {
        skipCert = true
    }
    
    trojan := models.Trojan{
        Type: "trojan",
        Tag: tag,
        Password: password,
        Server: serverUrl.Hostname(),
        Server_port: port,
        Tls: &models.Tls{
            Enabled: true,
            Insecure: skipCert,
            Server_name: params.Get("sni"),
        },
    }
    
    trojanMap, err := Struct2map(trojan, "trojan")
    if err != nil {
        
        return nil, err
    }
    
    return trojanMap, nil
}