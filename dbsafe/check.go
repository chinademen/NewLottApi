package dbsafe

import (
	"common"
	_ "fmt"
	"strings"
)

/**
 * 数据加密
 */
func Encode(sStr string) string {
	return Rsa.RsaEncryptString(sStr)
}

/**
 * 数据解密
 */
func Decode(sStr string) string {
	sStr = strings.Replace(sStr, "%2B", "+", -1)
	return Rsa.RsaDecryptString(sStr, "PKCS8")
}

/**
 * GetId
 */
func GetId() string {
	return common.GetKeyId()
}
