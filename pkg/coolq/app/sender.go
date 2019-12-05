package app

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"github.com/wppzxc/taobao_links/pkg/coolq/types"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	fileFormat = "20060102150405"
)

func CoolQMessageSend(msg types.Message, users []string) error {
	isImage := isImageMessage(msg.Message)
	if isImage {
		imageUrl := getImageUrl(msg.Message)
		tmpfile, err := SaveImage(imageUrl)
		if err != nil {
			return fmt.Errorf("Error in create tmp image file : %s ", err)
		}
		defer os.Remove(tmpfile.Name())
		errMsg := ""
		// send image
		for _, u := range users {
			if len(u) == 0 {
				continue
			}
			err := SendImage(tmpfile, u)
			if err != nil {
				errMsg = errMsg + " : " + err.Error()
			}
		}
		if len(errMsg) == 0 {
			return nil
		}
		return fmt.Errorf("Error in send image to users : %s ", errMsg)
	} else {
		errMsg := ""
		// send message
		for _, u := range users {
			if len(u) == 0 {
				continue
			}
			if err := SendMessage(msg.Message, u); err != nil {
				errMsg = errMsg + " : " + err.Error()
			}
		}
		if len(errMsg) == 0 {
			return nil
		}
		return fmt.Errorf("Error in send message to users : %s ", errMsg)
	}
	return nil
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
