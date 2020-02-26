package wechat

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/wppzxc/taobao_links/pkg/features/wechat/app"
	"github.com/wppzxc/taobao_links/pkg/features/wechat/types"
	"net/url"
	"strconv"
	"strings"
)

type Wechat struct {
	ParentWindow    *walk.MainWindow
	MainPage        *TabPage
	WebSocketUrl    *walk.LineEdit
	TaoKouLingTitle *walk.LineEdit
	Groups          *walk.TextEdit
	Users           *walk.TextEdit
	// 屏蔽关键词
	FilterNo        *walk.LineEdit
	// 筛选关键词
	FilterYes       *walk.LineEdit
	AutoImport      *walk.PushButton
	MoveToLeft      *walk.PushButton
	MoveToRight     *walk.PushButton
	Start           *walk.PushButton
	Stop            *walk.PushButton
	SendInterval    *walk.LineEdit
	StopCh          chan struct{}
}

func GetWechatPage() *Wechat {
	wechat := &Wechat{
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
	wechat.MainPage = &TabPage{
		Title:  "微信转发",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: wechat,
			AutoSubmit: true,
		},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "server端 websocket 地址：",
							},
						},
					},
					LineEdit{
						AssignTo: &wechat.WebSocketUrl,
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "替换淘口令标题(不填写则不替换)：",
							},
							LineEdit{
								AssignTo: &wechat.TaoKouLingTitle,
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "屏蔽关键词(多个用 / 分隔)：",
							},
							LineEdit{
								AssignTo: &wechat.FilterNo,
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "筛选关键词(多个用 / 分隔)：",
							},
							LineEdit{
								AssignTo: &wechat.FilterYes,
							},
						},
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "接收微信号（多个用 / 分隔）：",
							},
							Composite{
								Layout: HBox{Margins: Margins{}},
								Children: []Widget{
									HSpacer{},
									TextLabel{
										Text: "发送间隔（ms）：",
									},
									LineEdit{
										AssignTo: &wechat.SendInterval,
										MaxSize:  Size{50, 0},
										Text:     "1000",
									},
								},
							},
						},
					},
					TextEdit{
						AssignTo: &wechat.Groups,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "转发列表（多个用户用 / 分隔）：",
							},
							PushButton{
								Text:      "自动获取",
								AssignTo:  &wechat.AutoImport,
								OnClicked: wechat.AutoImportUsers,
							},
							PushButton{
								Text:      "移到左上角",
								AssignTo:  &wechat.MoveToLeft,
								OnClicked: wechat.MoveToLeftTop,
							},
							PushButton{
								Text:      "移到右上角",
								AssignTo:  &wechat.MoveToRight,
								OnClicked: wechat.MoveToRightTop,
							},
						},
					},
					TextEdit{
						AssignTo: &wechat.Users,
						VScroll:  true,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text:      "开始",
						AssignTo:  &wechat.Start,
						OnClicked: wechat.StartWork,
					},
					PushButton{
						Text:      "停止",
						AssignTo:  &wechat.Stop,
						OnClicked: wechat.StopWork,
						Enabled:   false,
					},
				},
			},
		},
	}
	return wechat
}

func (w *Wechat) AutoImportUsers() {
	users := app.AutoImportUsers(true, true)
	w.Users.SetText(users)
	fmt.Println("auto import users")
}

func (w *Wechat) MoveToLeftTop() {
	users := w.GetUsers()
	if len(users) == 0 {
		errMsg := "未指定转发用户！"
		walk.MsgBox(w.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}
	for _, u := range users {
		app.MoveUserToLeftTop(u)
	}
}

func (w *Wechat) MoveToRightTop() {
	users := w.GetUsers()
	if len(users) == 0 {
		errMsg := "未指定转发用户！"
		walk.MsgBox(w.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}
	for _, u := range users {
		app.MoveUserToRightTop(u)
	}
}

func (w *Wechat) StopWork() {
	fmt.Println("stop wechat trans work")
	close(w.StopCh)
	w.SetUIEnable(true)
}

func (w *Wechat) GetGroups() []string {
	data := w.Groups.Text()
	if len(data) == 0 {
		return nil
	}
	groups := strings.Split(data, "/")
	return groups
}

func (w *Wechat) GetUsers() []string {
	data := w.Users.Text()
	if len(data) == 0 {
		return nil
	}
	users := strings.Split(data, "/")
	return users
}

func (w *Wechat) SetUIEnable(enable bool) {
	w.Groups.SetEnabled(enable)
	w.Users.SetEnabled(enable)
	w.Start.SetEnabled(enable)
	w.WebSocketUrl.SetEnabled(enable)
	w.AutoImport.SetEnabled(enable)
	// stop is unEnable with others
	w.Stop.SetEnabled(!enable)
}

func (w *Wechat) StartWork() {
	fmt.Println("start wechat work!")
	wsUrl := w.WebSocketUrl.Text()
	if len(wsUrl) == 0 {
	Login:
		for {
			errMsg := "未指定 websocket url, 是否已登录本地酷Q？ "
			result := walk.MsgBox(w.ParentWindow, "Warning", errMsg, walk.MsgBoxOKCancel)
			switch result {
			case win.IDOK:
				fmt.Println("click already login coolq")
				if ok := app.CheckCoolqLogined(); !ok {
					fmt.Println("coolq not logined, please retry! ")
					continue
				}
				wsUrl = "ws://127.0.0.1:6700/event/"
				w.WebSocketUrl.SetText(wsUrl)
				break Login
			case win.IDCANCEL:
				fmt.Println("cancel login coolq")
				return
			default:
				fmt.Println("close login coolq")
				return
			}
		}
	}

	u, err := url.Parse(wsUrl)
	if err != nil {
		walk.MsgBox(w.ParentWindow, "Error", err.Error(), walk.MsgBoxIconError)
		return
	}

	types.Host = u.Host
	types.Port = u.Port()

	tlkTitle := w.TaoKouLingTitle.Text()
	filterNo := strings.Split(w.FilterNo.Text(), "/")
	filterYes := strings.Split(w.FilterYes.Text(), "/")

	if len(filterNo) == 1 && filterNo[0] == "" {
		filterNo = nil
	}

	if len(filterYes) == 1 && filterYes[0] == "" {
		filterYes = nil
	}

	intervalStr := w.SendInterval.Text()
	interval, _ := strconv.Atoi(intervalStr)
	if interval == 0 {
		fmt.Println("sendInterval can't be 0, reset to 1000")
		interval = 1000
		w.SendInterval.SetText("1000")
	}

	groups := w.GetGroups()
	if len(groups) == 0 {
		errMsg := "未指定QQ群ID！"
		walk.MsgBox(w.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}
	users := w.GetUsers()
	if len(users) == 0 {
		errMsg := "未指定转发用户！"
		walk.MsgBox(w.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}
	w.StopCh = make(chan struct{})
	w.SetUIEnable(false)
	if err := app.Start(wsUrl, groups, users, interval, tlkTitle, filterNo, filterYes, w.StopCh); err != nil {
		walk.MsgBox(w.ParentWindow, "Error", err.Error(), walk.MsgBoxIconError)
		w.SetUIEnable(true)
	}
	fmt.Println("wechat trans work started!")
}

func (w *Wechat) GetMainPage() *TabPage {
	return w.MainPage
}
