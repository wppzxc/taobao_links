package utils

import "strings"

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
