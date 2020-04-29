package dataokeapi

import (
	"crypto/md5"
	"encoding/hex"
)

func MakeSign(param string, appSecret string) string {
	str := param + "&key=" + appSecret
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
