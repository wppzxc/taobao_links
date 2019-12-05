package test

import (
	"fmt"
	"github.com/lxn/win"
	"github.com/shirou/w32"
	"strings"
	"syscall"
	"testing"
)

const (
	wxClass = "ChatWnd"
	// qq must remove the "QQ"
	qqClass = "TXGuiFoundation"
)

var Users []string

func TestGetWXUsers(t *testing.T) {
	mainH := win.GetDesktopWindow()
	ok := win.EnumChildWindows(mainH, Callback(EnumMainTVWindow_cn), 0)
	fmt.Println("EnumChildWindows ok : ", ok)
	fmt.Println("users are : ", Users)
}

func TestSplit(t *testing.T) {
	str := "123/42342/dasda"
	strs := strings.Split(str, "/")
	fmt.Println(strs)
}

func Callback(fn interface{}) uintptr {
	return syscall.NewCallback(fn)
}

func EnumMainTVWindow_cn(hWnd w32.HWND, UsersParam uintptr) uintptr {
	n := make([]uint16, 256)
	p := &n[0]
	name := w32.GetWindowText(hWnd)
	_, err := win.GetClassName(win.HWND(hWnd), p, 255)
	if err != nil {
		fmt.Println("Error in get class : ", err)
	}
	class := syscall.UTF16ToString(n)
	if class == wxClass {
		addUser(name)
	}
	return 1
}

func addUser(user string) {
	Users = append(Users, user)
}
