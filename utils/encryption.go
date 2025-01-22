package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncryptionMd5(str string) (string, error) {
	// 初始化MD5哈希器
	h := md5.New()
	// 将输入字符串写入哈希器
	_, err := h.Write([]byte(str))
	if err != nil {
		// 如果写入过程中出现错误，返回空字符串和错误
		return "", err
	}
	// 计算哈希值，并将其作为十六进制字符串返回
	return hex.EncodeToString(h.Sum(nil)), nil
}