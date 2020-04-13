package utils

import (
	"regexp"
	"strings"
)

func ParseData(data string) map[string]string {
	kvs := strings.Split(data, ";")
	cookies := make(map[string]string)
	for _, kv := range kvs {
		kv = strings.Trim(kv, " ")
		strs := strings.Split(kv, "=")
		cookies[strs[0]] = strs[1]
	}
	return cookies
}

func TranMoneySep(str string) string {
	reg := regexp.MustCompile(`[￥$][a-zA-Z0-9]{11}[￥$]`)
	oldStr := reg.FindString(str)
	var newStr string
	if strings.Contains(oldStr, "￥") {
		newStr = "(" + oldStr[3:len(oldStr)-3] + ")"
	} else {
		newStr = "(" + oldStr[1:len(oldStr)-1] + ")"
	}
	result := strings.Replace(str, oldStr, newStr, 1)
	return result
}
