package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/lxn/win"
)

var msg = "[图片=C:\\Users\\Administrator\\Desktop\\全新版本2019.1.30.1\\全新版本2019.1.30.1\\Cache\\Image\\670650fa-2217-4914-9b0c-1c6c3bfa56de.jpg]\r"

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
	desktopH := win.GetDesktopWindow()
	desktopRect := getWindowRect(desktopH)
	if desktopRect == nil {
		fmt.Println("Error in get desktop rect")
		return
	}

	hwnd := getWindowHWND(user)
	rect := getWindowRect(hwnd)
	if rect == nil {
		return
	}
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top
	moveRightX := desktopRect.Right - 1 - width
	//MoveUserToLeftTop(user)
	win.MoveWindow(hwnd, moveRightX, 0, width, height, false)
}

func TestExecuteMsg(t *testing.T) {
	msg := "**&^7788123test哈哈[CQ:image,file=9F3E0555B2FD9D2AD84539BEC47D8EA8.jpg,url=https://c2cpicdw.qpic.cn/offpic_new/980726589//72d93c83-69e5-497f-a3a4-6c5af5ed9f5f/0?vuin=1094187700&amp;term=2][CQ:image,file=9F3E0555B2FD9D2AD84539BEC47D8EA8.jpg,url=https://c2cpicdw.qpic.cn/offpic_new/980726589//72d93c83-69e5-497f-a3a4-6c5af5ed9f5f/0?vuin=1094187700&amp;term=2]http://www.baidu.com 快来啊123123"
	msgs := executeMessage(msg)
	fmt.Printf("msgs are : %#v", msgs)
}

func TestDecodeTaoKouLing(t *testing.T) {
	var tlk = "(BlqU121dKhp)"
	resp, err := http.Post(NianchuApiDecodeTaoKouLing,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("password_content=%s", tlk)))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in get resp data")
	}
	ncdresp := &NianchuDecodeResp{}
	if err := json.Unmarshal(data, ncdresp); err != nil {
		fmt.Println("Error in decode resp data")
	}
	fmt.Println(ncdresp.Data.Url)
}

func TestGenerateTaoKouLing(t *testing.T) {
	var url = "https://uland.taobao.com/taolijin/edetail?eh=ff6%2BAJTbF9WZuQF0XRz0iAXoB%2BDaBK5LQS0Flu%2FfbSp4QsdWMikAalrisGmre1Id0BFAqRODu11c9Nu1f984ALaUIeP3d6e2rwGSVh%2BekNlwjekHceFP2joQcqlJkEYSq25ItmKg62%2FQEQGZ%2BGEk9oEiRQvChEsh1QAU0Wz34ogNPFgv0tpQ0Q53%2B%2BKkcCZn3okcKRexRJpPpkxQH8hSJJv5JU%2ByyOgUNBWkg0gYTszgyVWw6NG5lMaJ3eKR%2F0jAcf6XM9SJXioUaXUnsIA1uo2nmxfVvdVf8UWdNg6VcHrluAYBRglsbQ%3D%3D&union_lens=lensId%3A0b01c26a_0c76_16fa908feed_dafe%3Btraffic_flag%3Dlm&activityId=361d30b13aab477e98c253fed24a0674&un=8a050443537b55c34c6e4bed9157ce66&share_crt_v=1&ut_sk=1.utdid_24551056_1579088939323.TaoPassword-Outside.taoketop&spm=a211b4.24665425&sp_tk=77+leTk3czEyWU15UVnvv6U=&visa=13a09278fde22a2e&disablePopup=true&disableSJ=1"
	var title = "test_title"
	resp, err := http.Post(NianchuApiGenerateTaoKouLing,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("url=%s&text=%s", url, title)))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in get resp data")
	}
	ncdresp := &NianchuGenerateResp{}
	if err := json.Unmarshal(data, ncdresp); err != nil {
		fmt.Println("Error in decode resp data")
	}
	tlk := strings.Replace(ncdresp.Data.Data.Model, "￥", "(", 1)
	tlk = strings.Replace(tlk, "￥", ")", 1)
	fmt.Println(tlk)
}
func TestChangeTaoKouLing(t *testing.T) {
	var tklTitle = "闹闹福利专享"
	var msg = `0.9！10支铅笔+5块橡皮
+1支水彩颜料
(yvyv12Ysdi4)`

	reg := regexp.MustCompile(`[(][a-zA-Z0-9]{11}[)]`)
	oldTkl := reg.FindString(msg)
	link := getTklLink(oldTkl)
	if len(link) == 0 {
		fmt.Println("Error in change taokouling: Error in get old taokouling link ")
	}
	newTkl := generateTklByTitle(link, tklTitle)
	if len(newTkl) == 0 {
		fmt.Println("Error in change taokouling: Error in get new taokouling code ")
	}
	newMsg := strings.Replace(msg, oldTkl, newTkl, 1)
	println(newMsg)
}
