package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wppzxc/taobao_links/pkg/features/PDDUserNumber"
	"github.com/wppzxc/taobao_links/pkg/features/coolq"
	"github.com/wppzxc/taobao_links/pkg/features/dataoke"
	"github.com/wppzxc/taobao_links/pkg/features/duoduojinbao"
	"github.com/wppzxc/taobao_links/pkg/features/goodsSearch"
	"github.com/wppzxc/taobao_links/pkg/features/haodanku"
	"github.com/wppzxc/taobao_links/pkg/features/taokeyi"
	"github.com/wppzxc/taobao_links/pkg/features/taokouling"
	"github.com/wppzxc/taobao_links/pkg/features/yituike"
	"github.com/wppzxc/taobao_links/pkg/version"
	"time"
)

type Feature interface {
	GetMainPage() *TabPage
}

func main() {
	mw := &walk.MainWindow{}
	now := time.Now().Unix()
	license := &License{}
	for {
		license = GetLocalLicense()
		if license == nil {
			fmt.Println("local license is invalid! ")
			licenseDlg := &LicenseDialog{
				MainWM: &walk.MainWindow{},
				License: &walk.TextEdit{},
				ActiveBtn: &walk.PushButton{},
			}
			if _, err := (MainWindow{
				Title:    "激活",
				Size:   Size{400, 200},
				AssignTo: &licenseDlg.MainWM,
				Layout:   VBox{},
				Children: []Widget{
					TextLabel{
						Text: "激活码......",
					},
					TextEdit{
						AssignTo: &licenseDlg.License,
					},
					PushButton{
						Text: "激活",
						AssignTo: &licenseDlg.ActiveBtn,
						OnClicked: licenseDlg.Active,
					},
				},
			}).Run(); err != nil {
				fmt.Println(err)
			}
			continue
		}
		if license.ExpireTimestamp > now {
			fmt.Println("local license ok !")
			break
		} else {
			fmt.Println("local license is expired !")
			if err := ResetLocalLicense(); err != nil {
				fmt.Println("Error in reset local license !", err)
			}
		}
	}


	// init dataoke
	dtk := dataoke.GetDataokePage()
	// init haodanku
	hdk := haodanku.GetHaodankuPage()
	// init duoduojinbao
	ddjb := duoduojinbao.GetDuoduojinbaoPage()
	// init taokeyi
	tky := taokeyi.GetTaokeyiPage()
	// init goodsSearch
	gs := goodsSearch.GetGoodsSearchPage()
	// init yituike
	ytk := yituike.GetYituikePage()
	if ok := ytk.LoadConfig(); ok {
		fmt.Println("load config from localFile config.yaml")
	}
	// init pddUserNumber
	pdun := PDDUserNumber.GetPDDUserNumberPage()
	// init taokouling
	tkl := taokouling.GetTaokoulingPage()
	// init coolq
	cq := coolq.GetCoolQPage()

	// bind mainWindow
	gs.ParentWindow = mw
	ytk.ParentWindow = mw
	pdun.ParentWindow = mw
	cq.ParentWindow = mw

	// all features
	featureMap := map[string]Feature{
		"dataoke": dtk,
		"haodanku": hdk,
		"duoduojinbao": ddjb,
		"taokeyi": tky,
		"goodSearch": gs,
		"yituike": ytk,
		"pddUserNumber": pdun,
		"taokouling": tkl,
		"coolq": cq,
	}

	// user features
	tabPages := []TabPage{}
	for _, f := range license.Feature {
		feature, ok := featureMap[f]
		if ok {
			tabPages = append(tabPages, *feature.GetMainPage())
		}
	}
	timeLayout := "2006-01-02"
	expireDate := time.Unix(license.ExpireTimestamp, 0).Format(timeLayout)

	if _, err := (MainWindow{
		Title:    getMainTitle() + expireDate,
		AssignTo: &mw,
		//Icon: "./assets/img/icon.ico",
		Size:   Size{700, 700},
		Layout: VBox{},
		Children: []Widget{
			TabWidget{
				Pages: tabPages,
			},
		},
	}).Run(); err != nil {
		fmt.Println(err)
	}
}

type LicenseDialog struct {
	MainWM    *walk.MainWindow
	License   *walk.TextEdit
	ActiveBtn *walk.PushButton
}

func (l *LicenseDialog) Active() {
	license := l.License.Text()
	if len(license) == 0 {
		walk.MsgBox(l.MainWM, "Error", fmt.Errorf("请填写激活码！").Error(), walk.MsgBoxIconError)
		return
	}
	if err := CheckLicense(license); err != nil {
		walk.MsgBox(l.MainWM, "Error", err.Error(), walk.MsgBoxIconError)
		return
	}
	fmt.Printf("active success! : %s\n", license)
	err := l.MainWM.Close()
	if err != nil {
		fmt.Println("close licenseDlg error : ", err)
	}
}

func getMainTitle() string {
	v := version.Get().String()
	return "淘宝客工具 " + v + " - license有效期至："
}
