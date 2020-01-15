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
	"github.com/wppzxc/taobao_links/pkg/license"
	"github.com/wppzxc/taobao_links/pkg/types"
	"github.com/wppzxc/taobao_links/pkg/utils/runtime"
	"github.com/wppzxc/taobao_links/pkg/version"
	"os"
	"time"
)

type Feature interface {
	GetMainPage() *TabPage
}

func main() {
	mw := &walk.MainWindow{}
	defer runtime.HandleCrash(mw)
	lic := &types.License{}
	lic = license.GetLocalLicense()
	if lic == nil {
		fmt.Println("local license is invalid! ")
		var err error
		lic, err = ShowLicenseDialog()
		if err != nil {
			fmt.Println("Error in active license : ", err)
		}
	}
	if err := license.CheckLicense(lic); err != nil {
		fmt.Printf("check license error %s \n", err)
		os.Exit(0)
	}
	fmt.Printf("license checked ok! : %#v \n", lic)

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
		"dataoke":       dtk,
		"haodanku":      hdk,
		"duoduojinbao":  ddjb,
		"taokeyi":       tky,
		"goodSearch":    gs,
		"yituike":       ytk,
		"pddUserNumber": pdun,
		"taokouling":    tkl,
		"coolq":         cq,
	}

	// user features
	tabPages := []TabPage{}
	for _, f := range lic.Feature {
		feature, ok := featureMap[f]
		if ok {
			tabPages = append(tabPages, *feature.GetMainPage())
		}
	}
	timeLayout := "2006-01-02"
	expireDate := time.Unix(lic.ExpireTimestamp, 0).Format(timeLayout)

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
	MainWM      *walk.MainWindow
	LicenseEdit *walk.TextEdit
	License     *types.License
	ActiveBtn   *walk.PushButton
}

func ShowLicenseDialog() (*types.License, error) {
	licenseDlg := &LicenseDialog{
		MainWM:      &walk.MainWindow{},
		License:     &types.License{},
		LicenseEdit: &walk.TextEdit{},
		ActiveBtn:   &walk.PushButton{},
	}
	if _, err := (MainWindow{
		Title:    "激活",
		Size:     Size{400, 200},
		AssignTo: &licenseDlg.MainWM,
		Layout:   VBox{},
		Children: []Widget{
			TextLabel{
				Text: "请输入激活码：",
			},
			TextEdit{
				AssignTo: &licenseDlg.LicenseEdit,
			},
			PushButton{
				Text:      "激活",
				AssignTo:  &licenseDlg.ActiveBtn,
				OnClicked: licenseDlg.Active,
			},
		},
	}).Run(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return licenseDlg.License, nil
}

func (l *LicenseDialog) Active() {
	lic := l.LicenseEdit.Text()
	if len(lic) == 0 {
		walk.MsgBox(l.MainWM, "Error", fmt.Errorf("请填写激活码！").Error(), walk.MsgBoxIconError)
		return
	}
	lice, err := license.CheckEncodeLicense(lic)
	if err != nil {
		walk.MsgBox(l.MainWM, "Error", err.Error(), walk.MsgBoxIconError)
		return
	}
	l.License = lice
	fmt.Printf("active success! : %s\n", lic)
	l.MainWM.Close()
	//l.MainWM.Closing()
}

func getMainTitle() string {
	v := version.Get().String()
	return "淘宝客工具 " + v + " - license有效期至："
}
