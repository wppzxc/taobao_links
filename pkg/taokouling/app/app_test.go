package app

import (
	"fmt"
	"testing"
)

func TestGetKouling(t *testing.T) {
	link := "http://t.cn/AigfcZbL"
	kl := GetKouling(link)
	fmt.Println(kl)
}
