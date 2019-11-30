package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wpp/taobao_links/pkg/PDDUserNumber"
	"github.com/wpp/taobao_links/pkg/dataoke"
	"github.com/wpp/taobao_links/pkg/duoduojinbao"
	"github.com/wpp/taobao_links/pkg/goodsSearch"
	"github.com/wpp/taobao_links/pkg/haodanku"
	"github.com/wpp/taobao_links/pkg/taokeyi"
	"github.com/wpp/taobao_links/pkg/taokouling"
	"github.com/wpp/taobao_links/pkg/yituike"
)

func main() {
	mw := &walk.MainWindow{}
	
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
	
	// bind mainWindow
	gs.ParentWindow = mw
	ytk.ParentWindow = mw
	pdun.ParentWindow = mw
	
	if _, err := (MainWindow{
		Title:    "淘宝客工具 v1.0.3",
		AssignTo: &mw,
		//Icon: "./assets/img/icon.ico",
		Size:    Size{600, 700},
		MaxSize: Size{700, 800},
		Layout:  VBox{},
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					*dtk.MainPage,
					*hdk.MainPage,
					*ddjb.MainPage,
					*tky.MainPage,
					*gs.MainPage,
					*ytk.MainPage,
					*pdun.MainPage,
					*tkl.MainPage,
				},
			},
		},
	}).Run(); err != nil {
		fmt.Println(err)
	}
}
