package control

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"sifu-box/ent"
	"sifu-box/model"
	"sifu-box/utils"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func FileList(key string, ent_client *ent.Client, logger *zap.Logger) ([]map[string]string, error) {

	templates, err := ent_client.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("查询模板失败: [%s]", err.Error()))
		return nil, err
	}
	template_list := []map[string]string{}
	for _, v := range templates {
		path := fmt.Sprintf("%x.json", md5.Sum([]byte(v.Name)))
		data := fmt.Sprintf(`%s%s%d`, v.Name, path, time.Now().Add(model.LINK_VALID_TIME).Unix())
		h := hmac.New(sha256.New, []byte(key))
		h.Write([]byte(data))
		signature := hex.EncodeToString(h.Sum(nil))
		template_list = append(template_list, map[string]string{"name": v.Name, "path": path, "signature": signature, "expire_time": fmt.Sprintf(`%d`, time.Now().Add(model.LINK_VALID_TIME).Unix())})
	}
	return template_list, nil
}
func FileDownload(work_dir, path, key, name, expire_time, signature string, logger *zap.Logger) ([]byte, error) {
	valid_time, err := strconv.ParseInt(expire_time, 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("时间转换失败: [%s]", err.Error()))
		return nil, err
	}

	if time.Now().Unix() > valid_time {
		return nil, fmt.Errorf("文件已过期")
	}
	data := fmt.Sprintf(`%s%s%d`, name, path, valid_time)
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	if hex.EncodeToString(h.Sum(nil)) != signature {
		return nil, fmt.Errorf("签名验证失败")
	}
	file_path := filepath.Join(work_dir, "sing-box", "config", path)
	content, err := utils.ReadFile(file_path)
	if err != nil {
		logger.Error(fmt.Sprintf("读取文件失败: [%s]", err.Error()))
		return nil, err
	}
	return content, nil
}
