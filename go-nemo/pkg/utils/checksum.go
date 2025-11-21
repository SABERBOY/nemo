package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
)

func GetMD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}

func GetCheckSum(nonce, curTime, appSecret string) string {
	str := appSecret + nonce + curTime
	h := sha1.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}
