package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func Sha256(secretKey string, str string) []byte {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(str))

	return h.Sum(nil)
}

func Base64Encode(str []byte) string {
	return base64.StdEncoding.EncodeToString(str)
}
