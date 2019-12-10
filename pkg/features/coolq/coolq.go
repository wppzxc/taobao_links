package coolq

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/wppzxc/taobao_links/pkg/features/coolq/app"
	"github.com/wppzxc/taobao_links/pkg/features/coolq/types"
	"strings"
)

type CoolQ struct {
	ParentWindow *walk.MainWindow
	MainPage     *TabPage
	WebSocketUrl *walk.LineEdit
	Groups       *walk.TextEdit
	Users        *walk.TextEdit
	LoginBtn     *walk.PushButton
	AutoImport   *walk.PushButton
	MoveTo       *walk.PushButton
	Start        *walk.PushButton
	Stop         *walk.PushButton
	StopCh       chan struct{}
}

func GetCoolQPage() *CoolQ {
	coolq := &CoolQ{
		WebSocketUrl: &walk.LineEdit{},
		Groups:       &walk.TextEdit{},
		Users:        &walk.TextEdit{},
		AutoImport:   &walk.PushButton{},
		Start:        &walk.PushButton{},
		Stop:         &walk.PushButton{},
	}
	coolq.MainPage = &TabPage{
		Title:  "酷Q消息转发",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: coolq,
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
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "接收群号（多个群组用 / 分隔）：",
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
								AssignTo:  &coolq.MoveTo,
								OnClicked: coolq.MoveToLeftTop,
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
	c.StopCh = make(chan struct{})
	c.SetUIEnable(false)
	if err := app.Start(wsUrl, groups, users, c.StopCh); err != nil {
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
