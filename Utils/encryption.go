package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncryptionMd5(str string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}