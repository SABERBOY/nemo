package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func CreateSudSignature(appId, appSecret, body, timestamp, nonce string) string {
	signContent := fmt.Sprintf("%s\n%s\n%s\n%s\n", appId, timestamp, nonce, body)
	mac := hmac.New(sha1.New, []byte(appSecret))
	mac.Write([]byte(signContent))
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifySudSignature(appId, appSecret, body, timestamp, nonce, signature string) bool {
	expectedSignature := CreateSudSignature(appId, appSecret, body, timestamp, nonce)
	return expectedSignature == signature
}
