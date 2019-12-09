package utils

import (
	"fmt"
	"github.com/lxn/win"
	"github.com/wppzxc/taobao_links/pkg/features/yituike/types"
	"strings"
	"time"
)

const (
	msgFormat = `%s

%s

开始时间: %s

免单网址：%s`
	timeFormat = "2006-01-02 15:04:05"
)

func GetMsg(preifx string, item *types.Item, link string) string {
	str := item.ExtendDocument
	str = strings.Replace(str, "#", "", -1)
	str = fmt.Sprintf(msgFormat, preifx, str, time.Unix(item.StartTime, 0).Format(timeFormat), link)
	return str
}

func CheckConfig(config *types.Config) bool {
	fmt.Println("test2")
	if len(config.Auth.Username) == 0 || len(config.Auth.Password) == 0 {
		return false
	}
	if !config.Fanli.Process.Start && !config.Fanli.Premonitor.Start {
		return false
	}
	if len(config.Receivers) == 0 {
		return false
	}
	fmt.Println("test3")
	return true
}

func GetDiffItems(oldItems []types.Item, newItems []types.Item) []types.Item {
	for _, o := range oldItems {
		for i := 0; i < len(newItems); i++ {
			if o.ExtendDocument == newItems[i].ExtendDocument && o.StartTime == newItems[i].StartTime {
				newItems = append(newItems[:i], newItems[i+1:]...)
				i--
				break
			}
		}
	}
	return newItems
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
