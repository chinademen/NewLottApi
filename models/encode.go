package models

import (
	"NewLottApi/dbsafe"
)

func EncodeString(sStr string) string {
	return dbsafe.Encode(sStr)
}

func DecodeString(sStr string) string {
	return dbsafe.Decode(sStr)
}
