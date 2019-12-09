package app

import (
	"fmt"
	"testing"
)

var msg = "[CQ:image,file=318000696E04B0529F892F66F19CEC9E.png,url=https://c2cpicdw.qpic.cn/offpic_new/980726589//cb6ae11d-fcd9-49d9-9631-099d856f08e9/0?vuin=1094187700&amp;term=2]"

func TestIsImage(t *testing.T) {
	image := isImageMessage(msg)
	fmt.Println("isImage : ", image)
}

func TestGetImageUrl(t *testing.T) {
	url := getImageUrl(msg)
	fmt.Println("get image url is : ", url)
}

func TestSaveImage(t *testing.T) {
	url := getImageUrl(msg)
	tmpFile, err := SaveImage(url)
	if err != nil {
		fmt.Println("Error in save image : ", err)
		return
	}
	fmt.Println(tmpFile.Name())
}

func TestMoveToLeftTop(t *testing.T) {
	user := "bigben"
	MoveUserToLeftTop(user)
}

func TestExecuteMsg(t *testing.T) {
	msg := "**&^7788123test哈哈[CQ:image,file=9F3E0555B2FD9D2AD84539BEC47D8EA8.jpg,url=https://c2cpicdw.qpic.cn/offpic_new/980726589//72d93c83-69e5-497f-a3a4-6c5af5ed9f5f/0?vuin=1094187700&amp;term=2][CQ:image,file=9F3E0555B2FD9D2AD84539BEC47D8EA8.jpg,url=https://c2cpicdw.qpic.cn/offpic_new/980726589//72d93c83-69e5-497f-a3a4-6c5af5ed9f5f/0?vuin=1094187700&amp;term=2]http://www.baidu.com 快来啊123123"
	msgs := executeMessage(msg)
	fmt.Printf("msgs are : %#v", msgs)
}