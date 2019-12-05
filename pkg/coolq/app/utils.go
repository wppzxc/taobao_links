package app

import (
	"fmt"
	"github.com/lxn/win"
	"github.com/shirou/w32"
	"strings"
	"syscall"
)

const (
	imageMsgPrefix = "CQ:image"
	wxClass        = "ChatWnd"
	// qq must remove the "QQ"
	qqClass = "TXGuiFoundation"
)

func isImageMessage(msg string) bool {
	if strings.Index(msg, imageMsgPrefix) >= 0 {
		return true
	}
	return false
}

func getImageUrl(msg string) string {
	strs := strings.Split(msg, "url=")
	if len(strs) >= 2 {
		return strs[1][:len(strs[1])-1]
	}
	return ""
}

func AutoImportUsers(weixin bool, qq bool) string {
	users := []string{}
	if weixin {
		wxUsers := importWXUsers()
		fmt.Println("wx users are : ", wxUsers)
		users = append(users, wxUsers...)
	}
	if qq {
		qqUsers := importQQUsers()
		fmt.Println("qq users are : ", qqUsers)
		users = append(users, qqUsers...)
	}
	userStr := ""
	for _, u := range users {
		userStr = userStr + u + "/"
	}
	return userStr
}

func importWXUsers() []string {
	users := []string{}
	callback := func(hWnd w32.HWND, UsersParam uintptr) uintptr {
		n := make([]uint16, 256)
		p := &n[0]
		name := w32.GetWindowText(hWnd)
		_, err := win.GetClassName(win.HWND(hWnd), p, 255)
		if err != nil {
			fmt.Println("Error in get class : ", err)
		}
		class := syscall.UTF16ToString(n)
		if class == wxClass && len(name) != 0 {
			users = append(users, name)
		}
		return 1
	}
	mainH := win.GetDesktopWindow()
	_ = win.EnumChildWindows(mainH, makeCallback(callback), 0)
	return users
}

func importQQUsers() []string {
	users := []string{}
	callback := func(hWnd w32.HWND, UsersParam uintptr) uintptr {
		n := make([]uint16, 256)
		p := &n[0]
		name := w32.GetWindowText(hWnd)
		_, err := win.GetClassName(win.HWND(hWnd), p, 255)
		if err != nil {
			fmt.Println("Error in get class : ", err)
		}
		class := syscall.UTF16ToString(n)
		if class == qqClass && len(name) != 0 && name != "QQ" && name != "TXMenuWindow" {
			users = append(users, name)
		}
		return 1
	}
	mainH := win.GetDesktopWindow()
	_ = win.EnumChildWindows(mainH, makeCallback(callback), 0)
	return users
}

func makeCallback(fn interface{}) uintptr {
	return syscall.NewCallback(fn)
}
