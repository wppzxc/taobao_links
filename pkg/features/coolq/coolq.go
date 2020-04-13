package coolq

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/wppzxc/taobao_links/pkg/features/coolq/app"
	"github.com/wppzxc/taobao_links/pkg/features/coolq/types"
	"strconv"
	"strings"
)

type CoolQ struct {
	ParentWindow      *walk.MainWindow
	MainPage          *TabPage
	WebSocketUrl      *walk.LineEdit
	TaoKouLingTitle   *walk.LineEdit
	Groups            *walk.TextEdit
	Users             *walk.TextEdit
	LoginBtn          *walk.PushButton
	AutoImport        *walk.PushButton
	MoveToLeft        *walk.PushButton
	MoveToRight       *walk.PushButton
	Start             *walk.PushButton
	Stop              *walk.PushButton
	SendInterval      *walk.LineEdit
	TranMoneySep      bool
	StringReplaceFrom string
	StringReplaceTo   string
	StopCh            chan struct{}
}

func GetCoolQPage() *CoolQ {
	coolq := &CoolQ{
		WebSocketUrl:    &walk.LineEdit{},
		TaoKouLingTitle: &walk.LineEdit{},
		Groups:          &walk.TextEdit{},
		Users:           &walk.TextEdit{},
		AutoImport:      &walk.PushButton{},
		Start:           &walk.PushButton{},
		Stop:            &walk.PushButton{},
		SendInterval:    &walk.LineEdit{},
	}
	coolq.MainPage = &TabPage{
		Title:  "酷Q消息转发",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: coolq,
			AutoSubmit: true,
			OnSubmitted: func() {
				fmt.Printf("%+v\n", coolq)
			},
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
								Text: "酷Q websocket 地址（默认使用本地）：",
							},
							PushButton{
								AssignTo: &coolq.LoginBtn,
								OnClicked: func() {
									if err := app.StartLocalCoolQ(); err != nil {
										walk.MsgBox(coolq.ParentWindow, "Warning", fmt.Sprintf("启动酷Q失败！%s \n", err), walk.MsgBoxIconWarning)
									}
								},
								Text: "启动酷Q",
							},
						},
					},
					LineEdit{
						AssignTo: &coolq.WebSocketUrl,
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "替换淘口令标题(不填写则不替换)：",
							},
							LineEdit{
								AssignTo: &coolq.TaoKouLingTitle,
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							CheckBox{
								Text: "启用￥、$ 转( )",
								Checked: Bind("TranMoneySep"),
							},
							HSpacer{},
							TextLabel{
								Text: "发送间隔（ms）：",
							},
							LineEdit{
								AssignTo: &coolq.SendInterval,
								MaxSize:  Size{50, 0},
								Text:     "1000",
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							TextLabel{
								Text: "文案替换（多个使用 / 分隔）：",
							},
							LineEdit{
								MinSize: Size{60, 0},
								Text: Bind("StringReplaceFrom"),
							},
							TextLabel{
								Text: "替换为",
							},
							LineEdit{
								MinSize: Size{60, 0},
								Text: Bind("StringReplaceFrom"),
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
								Text: "接收群号（多个群组用 / 分隔）：",
							},
						},
					},
					TextEdit{
						AssignTo: &coolq.Groups,
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
								AssignTo:  &coolq.AutoImport,
								OnClicked: coolq.AutoImportUsers,
							},
							PushButton{
								Text:      "移到左上角",
								AssignTo:  &coolq.MoveToLeft,
								OnClicked: coolq.MoveToLeftTop,
							},
							PushButton{
								Text:      "移到右上角",
								AssignTo:  &coolq.MoveToRight,
								OnClicked: coolq.MoveToRightTop,
							},
						},
					},
					TextEdit{
						AssignTo: &coolq.Users,
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
						AssignTo:  &coolq.Start,
						OnClicked: coolq.StartWork,
					},
					PushButton{
						Text:      "停止",
						AssignTo:  &coolq.Stop,
						OnClicked: coolq.StopWork,
						Enabled:   false,
					},
				},
			},
		},
	}
	return coolq
}

func (c *CoolQ) AutoImportUsers() {
	users := app.AutoImportUsers(true, true)
	c.Users.SetText(users)
	fmt.Println("auto import users")
}

func (c *CoolQ) MoveToLeftTop() {
	users := c.GetUsers()
	if len(users) == 0 {
		errMsg := "未指定转发用户！"
		walk.MsgBox(c.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}
	for _, u := range users {
		app.MoveUserToLeftTop(u)
	}
}

func (c *CoolQ) MoveToRightTop() {
	users := c.GetUsers()
	if len(users) == 0 {
		errMsg := "未指定转发用户！"
		walk.MsgBox(c.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}
	for _, u := range users {
		app.MoveUserToRightTop(u)
	}
}

func (c *CoolQ) StartWork() {
	fmt.Println("start coolq work!")
	wsUrl := c.WebSocketUrl.Text()
	if len(wsUrl) == 0 {
	Login:
		for {
			errMsg := "未指定 websocket url, 是否已登录本地酷Q？ "
			result := walk.MsgBox(c.ParentWindow, "Warning", errMsg, walk.MsgBoxOKCancel)
			switch result {
			case win.IDOK:
				fmt.Println("click already login coolq")
				if ok := app.CheckCoolqLogined(); !ok {
					fmt.Println("coolq not logined, please retry! ")
					continue
				}
				wsUrl = types.DefaultWebSocketUrl
				c.WebSocketUrl.SetText(wsUrl)
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
	tlkTitle := c.TaoKouLingTitle.Text()

	intervalStr := c.SendInterval.Text()
	interval, _ := strconv.Atoi(intervalStr)
	if interval == 0 {
		fmt.Println("sendInterval can't be 0, reset to 1000")
		interval = 1000
		c.SendInterval.SetText("1000")
	}

	groups := c.GetGroups()
	if len(groups) == 0 {
		errMsg := "未指定QQ群ID！"
		walk.MsgBox(c.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}
	users := c.GetUsers()
	if len(users) == 0 {
		errMsg := "未指定转发用户！"
		walk.MsgBox(c.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}
	stringReplaceFrom := strings.Split(c.StringReplaceFrom, "/")
	stringReplaceTo := strings.Split(c.StringReplaceTo, "/")

	if len(stringReplaceFrom) != len(stringReplaceTo) {
		errMsg := "文案替换填写不正确，请检查后重试！"
		walk.MsgBox(c.ParentWindow, "Error", errMsg, walk.MsgBoxIconError)
		return
	}

	c.StopCh = make(chan struct{})
	c.SetUIEnable(false)
	if err := app.Start(wsUrl, groups, users, interval, tlkTitle, c.TranMoneySep, stringReplaceFrom, stringReplaceTo, c.StopCh); err != nil {
		walk.MsgBox(c.ParentWindow, "Error", err.Error(), walk.MsgBoxIconError)
		c.SetUIEnable(true)
	}
	fmt.Println("coolq trans work started!")
}

func (c *CoolQ) StopWork() {
	fmt.Println("stop coolq trans work")
	close(c.StopCh)
	c.SetUIEnable(true)
}

func (c *CoolQ) GetGroups() []string {
	data := c.Groups.Text()
	if len(data) == 0 {
		return nil
	}
	groups := strings.Split(data, "/")
	return groups
}

func (c *CoolQ) GetUsers() []string {
	data := c.Users.Text()
	if len(data) == 0 {
		return nil
	}
	users := strings.Split(data, "/")
	return users
}

func (c *CoolQ) SetUIEnable(enable bool) {
	c.Groups.SetEnabled(enable)
	c.Users.SetEnabled(enable)
	c.Start.SetEnabled(enable)
	c.WebSocketUrl.SetEnabled(enable)
	c.AutoImport.SetEnabled(enable)
	// stop is unEnable with others
	c.Stop.SetEnabled(!enable)
}

func (c *CoolQ) GetMainPage() *TabPage {
	return c.MainPage
}
