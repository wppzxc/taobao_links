package qiwechat

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type QiWechat struct {
	ParentWindow    *walk.MainWindow
	MainPage        *TabPage
	WebSocketUrl    *walk.LineEdit
	TaoKouLingTitle *walk.LineEdit
	Groups          *walk.TextEdit
	Users           *walk.TextEdit
	// 屏蔽关键词
	FilterNo *walk.LineEdit
	// 筛选关键词
	FilterYes    *walk.LineEdit
	AutoImport   *walk.PushButton
	MoveToLeft   *walk.PushButton
	MoveToRight  *walk.PushButton
	Start        *walk.PushButton
	Stop         *walk.PushButton
	SendInterval *walk.LineEdit
	StopCh       chan struct{}
}

const (
	INPUT_MOUSE          = 0
	MOUSEEVENTF_LEFTDOWN = 0x0002
	MOUSEEVENTF_LEFTUP   = 0x0004
)

type MOUSEINPUT struct {
	dx, dy      int32
	mouseData   uint32
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	Mi   MOUSEINPUT
}

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	procSendInput = user32.NewProc("SendInput")
)

func sendMouseClick(x, y int32) {
	// 设置鼠标位置
	win.SetCursorPos(x, y)
	var inputs [2]INPUT

	inputs[0].Type = INPUT_MOUSE
	inputs[0].Mi.dwFlags = MOUSEEVENTF_LEFTDOWN

	inputs[1].Type = INPUT_MOUSE
	inputs[1].Mi.dwFlags = MOUSEEVENTF_LEFTUP

	ret, _, err := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)
	if ret == 0 {
		fmt.Println("SendInput failed:", err)
	}
	fmt.Println("SendInput ok")
}

func GetWechatPage() *QiWechat {
	qiWechat := &QiWechat{
		WebSocketUrl:    &walk.LineEdit{},
		TaoKouLingTitle: &walk.LineEdit{},
		Groups:          &walk.TextEdit{},
		Users:           &walk.TextEdit{},
		FilterNo:        &walk.LineEdit{},
		FilterYes:       &walk.LineEdit{},
		AutoImport:      &walk.PushButton{},
		Start:           &walk.PushButton{},
		Stop:            &walk.PushButton{},
		SendInterval:    &walk.LineEdit{},
	}
	qiWechat.MainPage = &TabPage{
		Title:  "一秒万次",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: qiWechat,
			AutoSubmit: true,
		},
		Children: []Widget{
			// Composite{
			// 	Layout: VBox{},
			// 	Children: []Widget{
			// 		HSpacer{},
			// 		Composite{
			// 			Layout: HBox{},
			// 			Children: []Widget{
			// 				TextLabel{
			// 					Text: "server端 websocket 地址：",
			// 				},
			// 			},
			// 		},
			// 		LineEdit{
			// 			AssignTo: &wechat.WebSocketUrl,
			// 		},
			// 		Composite{
			// 			Layout: HBox{},
			// 			Children: []Widget{
			// 				TextLabel{
			// 					Text: "替换淘口令标题(不填写则不替换)：",
			// 				},
			// 				LineEdit{
			// 					AssignTo: &wechat.TaoKouLingTitle,
			// 				},
			// 			},
			// 		},
			// 		Composite{
			// 			Layout: HBox{},
			// 			Children: []Widget{
			// 				TextLabel{
			// 					Text: "屏蔽关键词(多个用 / 分隔)：",
			// 				},
			// 				LineEdit{
			// 					AssignTo: &wechat.FilterNo,
			// 				},
			// 			},
			// 		},
			// 		Composite{
			// 			Layout: HBox{},
			// 			Children: []Widget{
			// 				TextLabel{
			// 					Text: "筛选关键词(多个用 / 分隔)：",
			// 				},
			// 				LineEdit{
			// 					AssignTo: &wechat.FilterYes,
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			// Composite{
			// 	Layout: VBox{},
			// 	Children: []Widget{
			// 		HSpacer{},
			// 		Composite{
			// 			Layout: HBox{},
			// 			Children: []Widget{
			// 				TextLabel{
			// 					Text: "接收微信号（多个用 / 分隔）：",
			// 				},
			// 				Composite{
			// 					Layout: HBox{Margins: Margins{}},
			// 					Children: []Widget{
			// 						HSpacer{},
			// 						TextLabel{
			// 							Text: "发送间隔（ms）：",
			// 						},
			// 						LineEdit{
			// 							AssignTo: &wechat.SendInterval,
			// 							MaxSize:  Size{50, 0},
			// 							Text:     "1000",
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 		TextEdit{
			// 			AssignTo: &wechat.Groups,
			// 		},
			// 	},
			// },
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "坐标列表（多个用 / 分隔）：",
							},
							TextLabel{
								Text: "按 鼠标中键 添加坐标",
							},
							PushButton{
								Text:      "全部清除",
								AssignTo:  &qiWechat.AutoImport,
								OnClicked: qiWechat.RemoveAllPosition,
							},
							PushButton{
								Text:      "开始执行",
								AssignTo:  &qiWechat.MoveToLeft,
								OnClicked: qiWechat.StartWork,
							},
						},
					},
					TextEdit{
						AssignTo: &qiWechat.Users,
						VScroll:  true,
						ReadOnly: true,
					},
				},
			},
			// Composite{
			// 	Layout: HBox{},
			// 	Children: []Widget{
			// 		HSpacer{},
			// 		PushButton{
			// 			Text:      "开始",
			// 			AssignTo:  &qiWechat.Start,
			// 			OnClicked: qiWechat.StartWork,
			// 		},
			// 		PushButton{
			// 			Text:      "停止",
			// 			AssignTo:  &qiWechat.Stop,
			// 			OnClicked: qiWechat.StopWork,
			// 			Enabled:   false,
			// 		},
			// 	},
			// },
		},
	}
	// 监听 F6事件
	go qiWechat.listenMouseCenterHotkey()
	return qiWechat
}

// listenMouseCenterHotkey 监听 mouse center 点击事件
// 点击鼠标中键时，记录当前鼠标坐标
func (w *QiWechat) listenMouseCenterHotkey() {
	var lastPressed time.Time
	// 设置去抖动的时间间隔 1s
	debounceDuration := 1 * time.Second
	for {
		// 监听 mouse center 被按下（true 表示按下，false 表示松开）
		if robotgo.AddEvent("center") {
			fmt.Println("mouse center detected, add position...")
			// 获取当前时间
			now := time.Now()
			// 如果上次按键与当前时间差小于 debounceDuration，则跳过此次按键事件
			if now.Sub(lastPressed) < debounceDuration {
				continue
			}
			// 更新上次按键时间
			lastPressed = now

			var pt win.POINT
			win.GetCursorPos(&pt)
			allPositions := w.Users.Text()
			newAllPositions := path.Join(allPositions, fmt.Sprintf("%d,%d", pt.X, pt.Y))
			w.Users.SetText(newAllPositions)
		}
	}
}

// RemoveAllPosition 移除所有坐标
func (w *QiWechat) RemoveAllPosition() {
	w.Users.SetText("")
	fmt.Println("remove all position")
}

// StartWork 开始执行, 将记录的坐标全都点击一遍
func (w *QiWechat) StartWork() {
	positionStr := w.Users.Text()
	ps := strings.Split(positionStr, "/")
	for _, p := range ps {
		pos := strings.Split(p, ",")
		x, _ := strconv.Atoi(pos[0])
		y, _ := strconv.Atoi(pos[1])
		sendMouseClick(int32(x), int32(y))
		// time.Sleep(5 * time.Millisecond)
	}
}

func (w *QiWechat) StopWork() {
	fmt.Println("stop qiWechat trans work")
	close(w.StopCh)
	w.SetUIEnable(true)
}

func (w *QiWechat) GetGroups() []string {
	data := w.Groups.Text()
	if len(data) == 0 {
		return nil
	}
	groups := strings.Split(data, "/")
	return groups
}

func (w *QiWechat) GetUsers() []string {
	data := w.Users.Text()
	if len(data) == 0 {
		return nil
	}
	users := strings.Split(data, "/")
	return users
}

func (w *QiWechat) SetUIEnable(enable bool) {
	w.Groups.SetEnabled(enable)
	w.Users.SetEnabled(enable)
	w.Start.SetEnabled(enable)
	w.WebSocketUrl.SetEnabled(enable)
	w.AutoImport.SetEnabled(enable)
	// stop is unEnable with others
	w.Stop.SetEnabled(!enable)
}

func (w *QiWechat) GetMainPage() *TabPage {
	return w.MainPage
}
