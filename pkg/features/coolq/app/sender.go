package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"github.com/wppzxc/taobao_links/pkg/features/coolq/types"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
)

const (
	fileFormat = "20060102150405"
)

const (
	NianchuApiDecodeTaoKouLing   = "http://admin.nianchu.net/api/5bc055393b4b1"
	NianchuApiGenerateTaoKouLing = "http://admin.nianchu.net/api/5bc044955940d"
)

func CoolQMessageSend(msg types.Message, users []string, interval int, tklTitle string) error {
	var errMsg = ""
	var sendMsgs []string
	msgs := executeMessage(msg.Message)
	if len(msgs) == 0 {
		return nil
	}
	// change taokouling title
	for _, m := range msgs {
		isImage := isImageMessage(m)
		if isImage {
			sendMsgs = append(sendMsgs, m)
			continue
		}
		newMsg := changeTaoKouLing(m, tklTitle)
		sendMsgs = append(sendMsgs, newMsg)
	}

	for _, u := range users {
		for _, m := range sendMsgs {
			isImage := isImageMessage(m)
			// send image to user
			if isImage {
				imageUrl := getImageUrl(m)
				tmpfile, err := SaveImage(imageUrl)
				if err != nil {
					errMsg = errMsg + " : " + fmt.Sprintf("Error in create tmp image file : %s ", err)
					continue
				}
				func() {
					defer tmpfile.Close()
					defer os.Remove(tmpfile.Name())
					if err := SendImage(tmpfile, u); err != nil {
						errMsg = errMsg + " : " + err.Error()
					}
				}()
			} else {
				if err := SendMessage(m, u); err != nil {
					errMsg = errMsg + " : " + err.Error()
				}
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	if len(errMsg) == 0 {
		return nil
	}
	return fmt.Errorf("Error in send message to users : %s ", errMsg)
}

func SendMessage(msg string, user string) error {
	if err := clipboard.WriteAll(msg); err != nil {
		return fmt.Errorf("Error on write to clipboard : %s ", err)
	}

	if err := send(user); err != nil {
		return err
	}

	return nil
}

func SendImage(img *os.File, user string) error {
	if img == nil {
		return fmt.Errorf("Warning in send image file : input image file is nil ")
	}

	_, err := exec.Command("file2clip", img.Name()).CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error in save image to clipboard %s ", err)
	}

	if err := send(user); err != nil {
		return err
	}
	return nil
}

func send(user string) error {
	p, err := syscall.UTF16PtrFromString(user)
	if err != nil {
		return fmt.Errorf("Error in get user chat window : %s ", err)
	}
	h2 := win.FindWindow(nil, p)
	if h2 == 0 {
		return fmt.Errorf("Error in get user %s, user not found ", user)
	}
	var ok = false
	for i := 0; i < 10; i++ {
		ok = SetForegroundWindow(h2)
		if ok {
			break
		}
	}

	if ok {
		robotgo.KeyTap("v", "ctrl")
		robotgo.KeyTap("enter")
		//return nil
	} else {
		return fmt.Errorf("Error in set window foreground %s in 10 times ", user)
	}
	return nil
}

func SaveImage(url string) (*os.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error in download image %s ", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if len(body) == 0 {
		return nil, fmt.Errorf("Error in download image, image is null ")
	}

	tmpfile, err := os.Create("./coolq_image_" + time.Now().Format(fileFormat) + ".png")
	if err != nil {
		return nil, fmt.Errorf("Error in create temp file %s ", err)
	}
	defer tmpfile.Close()

	_, err = io.Copy(tmpfile, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("Error in save image to tmpfile %s ", err)
	}

	return tmpfile, nil
}

func SetForegroundWindow(hWnd win.HWND) bool {
	hForeWnd := win.GetForegroundWindow()
	dwCurID := win.GetCurrentThreadId()
	dcID := int32(dwCurID)
	dwForeID := win.GetWindowThreadProcessId(hForeWnd, nil)
	dfID := int32(dwForeID)
	win.AttachThreadInput(dcID, dfID, true)
	win.ShowWindow(hWnd, win.SW_SHOWNORMAL)
	win.SetWindowPos(hWnd, win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE)
	win.SetWindowPos(hWnd, win.HWND_NOTOPMOST, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE)
	ok := win.SetForegroundWindow(hWnd)
	win.AttachThreadInput(dcID, dfID, false)
	return ok
}

func executeMessage(msg string) []string {
	msgs := []string{}

	// match image
	reg := regexp.MustCompile(`[[]CQ:image[a-zA-Z0-9_\-,=.:;?/&]+]`)
	strs := reg.FindAllString(msg, -1)
	length := len(strs)
	// if don't have image, return msg
	if length == 0 {
		msgs = append(msgs, msg)
		return msgs
	}
	// get message before image
	tmpStrs := strings.Split(msg, strs[0])
	if len(tmpStrs[0]) > 0 {
		msgs = append(msgs, tmpStrs[0])
	}

	// add images
	msgs = append(msgs, strs...)

	// get message after image
	tmpStrs2 := strings.Split(msg, strs[length-1])
	if len(tmpStrs[len(tmpStrs2)-1]) > 0 {
		msgs = append(msgs, tmpStrs[len(tmpStrs2)-1])
	}
	return msgs
}

func changeTaoKouLing(msg string, tklTitle string) string {
	if len(tklTitle) == 0 {
		return msg
	}

	reg := regexp.MustCompile(`[(][a-zA-Z0-9]{11}[)]`)
	oldTkl := reg.FindString(msg)
	if len(oldTkl) > 0 {
		link := getTklLink(oldTkl)
		if len(link) == 0 {
			return msg
		}

		newTkl := generateTklByTitle(link, tklTitle)
		newMsg := strings.Replace(msg, oldTkl, newTkl, 1)
		return newMsg
	}else {
		return msg
	}
}

type NianchuResp struct {
	Code int    `json:"code"`
	Msg  string `json:"Msg"`
}

type NianchuDecodeData struct {
	Content string `json:"content"`
	Url     string `json:"url"`
}

type NianchuDecodeResp struct {
	NianchuResp
	Data NianchuDecodeData `json:"data"`
}

type NianchuGenerateModelData struct {
	Model string `json:"model"`
}

type NianchuGenerateData struct {
	Data NianchuGenerateModelData `json:"data"`
}

type NianchuGenerateResp struct {
	NianchuResp
	Data NianchuGenerateData `json:"data"`
}

func getTklLink(tlk string) string {
	resp, err := http.Post(NianchuApiDecodeTaoKouLing,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("password_content=%s", tlk)))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	ncdresp := &NianchuDecodeResp{}
	if err := json.Unmarshal(data, ncdresp); err != nil {
		return ""
	}
	return ncdresp.Data.Url
}

func generateTklByTitle(url, title string) string {
	resp, err := http.Post(NianchuApiGenerateTaoKouLing,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("url=%s&text=%s", url, title)))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	ncdresp := &NianchuGenerateResp{}
	if err := json.Unmarshal(data, ncdresp); err != nil {
		return ""
	}
	tlk := strings.Replace(ncdresp.Data.Data.Model, "￥", "(", 1)
	tlk = strings.Replace(tlk, "￥", ")", 1)
	return tlk
}
