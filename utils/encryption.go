package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// EncryptionMd5 使用MD5算法对字符串进行加密。
// 它接受一个字符串作为输入，并返回加密后的MD5字符串表示形式。
// 如果加密过程中出现错误，它将返回错误。
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