package app

import (
	"fmt"
	"testing"
)

func TestGetNameMemberByLink(t *testing.T) {
	url := "https://mobile.yangkeduo.com/goods2.html?goods_id=34563598420"
	key := ""
	data, num := GetNameMemberByLink(url, key)
	fmt.Println(num)
	fmt.Println(data)
}
