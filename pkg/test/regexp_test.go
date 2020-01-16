package test

import (
	"fmt"
	"regexp"
	"testing"
	"time"
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

func TestTaoKouLing(t *testing.T) {
	var text = `0.9！10支铅笔+5块橡皮
+1(支水彩颜料)
(yvyv12Ysdi6)`
	reg := regexp.MustCompile(`[(][a-zA-Z0-9]{11}[)]`)
	result := reg.FindString(text)
	fmt.Println(result)

}

func TestChan(t *testing.T) {
	start := time.Now()
	c := make(chan interface{})
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(4 * time.Second)
		close(c)
	}()

	go func() {
		time.Sleep(3 * time.Second)
		ch1 <- 3
	}()

		time.Sleep(2 * time.Second)
		ch2 <- 5
	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	case <-ch1:
		fmt.Printf("ch1 case...")
	case <-ch2:
		fmt.Printf("ch2 case...")
		//default:
		//	fmt.Printf("default go...")
	}
}
