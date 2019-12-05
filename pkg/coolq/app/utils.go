package app

import "strings"

const (
	imageMsgPrefix = "CQ:image"
)

func isImageMessage(msg string) bool {
	if strings.Index(msg, imageMsgPrefix) >= 0 {
		return true
	}
	return false
}

func getImageUrl(msg string) string {
	strs := strings.Split(msg, "url=")
	if len(strs) >= 2 {
		return strs[1][:len(strs[1])-1]
	}
	return ""
}
