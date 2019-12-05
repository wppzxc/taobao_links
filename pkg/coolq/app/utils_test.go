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