package test

import (
	"fmt"
	"regexp"
	"testing"
)

const (
	str = "[CQ:image,file=9F3E0555B2FD9D2AD84539BEC47D8EA8.jpg,url=https://c2cpicdw.qpic.cn/offpic_new/980726589//72d93c83-69e5-497f-a3a4-6c5af5ed9f5f/0?vuin=1094187700&amp;term=2]http://www.baidu.com 快来啊123123"
)

func TestRegExp(t *testing.T) {
	//match, _ := regexp.Match("[[]CQ:image[A-Z0-9,=.:;?/&]+]", []byte(str))
	reg := regexp.MustCompile(`[[]CQ:image[a-zA-Z0-9_\-,=.:;?/&]+]`)
	if reg == nil {
		fmt.Println("regexp error !")
		return
	}
	//match, _ := regexp.Match("[[]CQ:image[A-Z0-9,=.:;?/&]+", []byte(str))
	result := reg.FindAllString(str, -1)
	fmt.Println("match is : ", result)
}
