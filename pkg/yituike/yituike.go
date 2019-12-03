package yituike

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wppzxc/taobao_links/pkg/yituike/app"
	"github.com/wppzxc/taobao_links/pkg/yituike/filter"
	"github.com/wppzxc/taobao_links/pkg/yituike/premonitor"
	"github.com/wppzxc/taobao_links/pkg/yituike/process"
	"github.com/wppzxc/taobao_links/pkg/yituike/types"
	"github.com/wppzxc/taobao_links/pkg/yituike/utils"
	"io/ioutil"
	"os"
	"path"
	"sigs.k8s.io/yaml"
	"strconv"
)

const (
	configFile = "config.yaml"
)

type Yituike struct {
	ParentWindow    *walk.MainWindow
	MainPage        *TabPage
	RunBtn          *walk.PushButton
	Config          *types.Config
	RefreshInterval string
	SendInterval    string
	StopCh          chan struct{}
	FilterQuanNum   string
	UIFileName      *walk.LineEdit
	Running         bool
}

func GetYituikePage() *Yituike {
	yituike := &Yituike{
		Config: new(types.Config),
	}
	yituike.LoadConfig()
	yituike.MainPage = &TabPage{
		Title:  "易推客发单",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: yituike,
			AutoSubmit: true,
		},
		Children: []Widget{
			// auth
			Composite{
				MaxSize: Size{0, 400},
				Layout:  VBox{},
				Children: []Widget{
					Label{
						Text: "认证信息",
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Composite{
								Layout: HBox{},
								Children: []Widget{
									HSpacer{},
									Label{
										Text: "用户名: ",
									},
									LineEdit{
										Text: Bind("Config.Auth.Username"),
									},
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									HSpacer{},
									Label{
										Text: "密码: ",
									},
									LineEdit{
										Text: Bind("Config.Auth.Password"),
									},
								},
							},
						},
					},
				},
			},
			// receivers
			Composite{
				MaxSize: Size{0, 200},
				Layout:  VBox{},
				Children: []Widget{
					Label{
						Text: "接收者信息",
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							HSpacer{},
							PushButton{
								Text:      "选择文件",
								OnClicked: yituike.ChooseFile,
							},
							LineEdit{
								AssignTo: &yituike.UIFileName,
								ReadOnly: true,
								Text:     Bind("Config.ReceiversFile"),
							},
						},
					},
				},
			},
			// fanli
			Composite{
				MinSize: Size{0, 300},
				Layout:  VBox{},
				Children: []Widget{
					Label{
						Text: "返利信息",
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "优惠劵数量阈值：",
							},
							LineEdit{
								Text: Bind("FilterQuanNum"),
							},
							Label{
								Text: "刷新间隔：",
							},
							LineEdit{
								MinSize: Size{60, 0},
								Text:    Bind("RefreshInterval"),
							},
							Label{
								Text: "发送间隔：",
							},
							LineEdit{
								MinSize: Size{60, 0},
								Text:    Bind("SendInterval"),
							},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							Label{
								Text: "预告商品：",
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									CheckBox{
										Text:    "开启",
										Name:    "Pre",
										Checked: Bind("Config.Fanli.Premonitor.Start"),
									},
									Label{
										Text: "消息前缀",
									},
									LineEdit{
										Enabled: Bind("Pre.Checked"),
										MinSize: Size{60, 0},
										Text:    Bind("Config.Fanli.Premonitor.MsgPrefix"),
									},
								},
							},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							Label{
								Text: "进行中商品：",
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									CheckBox{
										Text:    "开启",
										Name:    "Pro",
										Checked: Bind("Config.Fanli.Process.Start"),
									},
									Label{
										Text: "消息前缀",
									},
									LineEdit{
										Enabled: Bind("Pro.Checked"),
										MinSize: Size{60, 0},
										Text:    Bind("Config.Fanli.Process.MsgPrefix"),
									},
								},
							},
						},
					},
				},
			},
			// start
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text:      "开始",
						OnClicked: yituike.Run,
						AssignTo:  &yituike.RunBtn,
					},
					PushButton{
						Text:      "停止",
						OnClicked: yituike.Stop,
					},
				},
			},
		},
	}
	return yituike
}

func (y *Yituike) ChooseFile() {
	dlg := new(walk.FileDialog)
	dlg.FilePath = y.Config.ReceiversFile
	dlg.Filter = "配置文件 (*.yaml)"
	dlg.Title = "选择文件"
	if ok, err := dlg.ShowOpen(y.ParentWindow); err != nil {
		fmt.Printf("Error in open file : %s\n", err)
	} else if !ok {
		return
	}

	y.Config.ReceiversFile = dlg.FilePath
	y.UIFileName.SetText(y.Config.ReceiversFile)
	fmt.Println(y.Config.ReceiversFile)
}

func (y *Yituike) Run() {
	conf := &types.Config{}
	// get Receivers
	data, _ := ioutil.ReadFile(y.Config.ReceiversFile)
	if err := yaml.Unmarshal(data, conf); err != nil {
		fmt.Println(err)
		return
	}
	y.Config.Receivers = conf.Receivers

	// get auth
	y.Config.Auth.Url = types.AuthUrl

	// get fanli
	filterQuanNum, _ := strconv.Atoi(y.FilterQuanNum)
	y.Config.Fanli.FilterQuanNum = filterQuanNum
	fmt.Printf("filterNum is %d", y.Config.Fanli.FilterQuanNum)
	ref, _ := strconv.ParseInt(y.RefreshInterval, 10, 64)
	sen, _ := strconv.ParseInt(y.SendInterval, 10, 64)
	y.Config.Fanli.RefreshInterval = ref
	y.Config.Fanli.SendInterval = sen
	y.Config.Fanli.Premonitor.Url = types.PreUrl
	y.Config.Fanli.Process.Url = types.ProUrl
	fmt.Println("test1")
	fmt.Printf("config is %#v", y.Config)
	if ok := utils.CheckConfig(y.Config); !ok {
		fmt.Println("test5")
		return
	}
	y.SaveConfig()
	fmt.Println("save config ok !")
	y.Running = true
	y.Start()
	fmt.Printf("config is %#v", y.Config)
}

func (y *Yituike) Start() {
	y.SetUIEnable(false)
	y.StopCh = make(chan struct{})
	fs := []func(){}
	if y.Config.Fanli.Process.Start {
		pro := process.Processer{
			Config: y.Config,
			Filter: filter.FilterByQuan,
		}
		fs = append(fs, pro.StartProcess)
	}
	if y.Config.Fanli.Premonitor.Start {
		pre := premonitor.Premonitor{
			Config: y.Config,
			Filter: filter.FilterByQuan,
		}
		fs = append(fs, pre.StartPremonitor)
	}
	app.Start(y.Config, fs, y.StopCh)
}

func (y *Yituike) Stop() {
	if y.Running == false {
		return
	}
	y.Running = false
	fmt.Println("stop !")
	y.SetUIEnable(true)
	close(y.StopCh)
}

func (y *Yituike) SetUIEnable(enable bool) {
	y.RunBtn.SetEnabled(enable)
	y.UIFileName.SetEnabled(enable)
}

func (y *Yituike) LoadConfig() bool {
	file := configFile
	_, err := os.Stat(file)
	if os.IsPermission(err) {
		fmt.Println(err)
		return false
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(file), 0755); err != nil {
			fmt.Println(err)
		}
		_, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
		}
		return false
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if len(data) == 0 {
		return false
	}
	conf := &types.Config{}
	if err := yaml.Unmarshal(data, conf); err != nil {
		fmt.Println(err)
		return false
	}
	y.Config = conf
	return true
}

func (y *Yituike) SaveConfig() {
	fmt.Println("start save config file...")
	file := configFile
	data, _ := yaml.Marshal(y.Config)
	ioutil.WriteFile(file, data, 0755)
	fmt.Println("end save config file...")
}
