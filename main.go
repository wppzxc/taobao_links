package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wpp/taobao_links/pkg/dataoke"
	"github.com/wpp/taobao_links/pkg/duoduojinbao"
	"github.com/wpp/taobao_links/pkg/goodsSearch"
	"github.com/wpp/taobao_links/pkg/haodanku"
	"github.com/wpp/taobao_links/pkg/taokeyi"
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
	
	// bind mainWindow
	gs.ParentWindow = mw
	ytk.ParentWindow = mw
	
	if _, err := (MainWindow{
		Title:    "淘宝客工具 v1.0.0",
		AssignTo: &mw,
		//Icon: "./assets/img/icon.ico",
		Size:    Size{350, 600},
		MaxSize: Size{350, 600},
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
				},
			},
		},
	}).Run(); err != nil {
		fmt.Println(err)
	}
}
