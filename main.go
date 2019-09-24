package main

import (
	"fmt"
	. "github.com/lxn/walk/declarative"
	"github.com/wpp/taobao_links/pkg/dataoke"
	"github.com/wpp/taobao_links/pkg/duoduojinbao"
	haodanku "github.com/wpp/taobao_links/pkg/haodanku"
	"github.com/wpp/taobao_links/pkg/taokeyi"
)

func main() {
	// init dataoke
	dtk := dataoke.GetDataokePage()
	// init haodanku
	hdk := haodanku.GetHaodankuPage()
	// init duoduojinbao
	ddjb := duoduojinbao.GetDuoduojinbaoPage()
	// init taokeyi
	tky := taokeyi.GetTaokeyiPage()
	if _, err := (MainWindow{
		Title: "商品链接获取工具",
		//Icon: "./assets/img/icon.ico",
		Size:   Size{350, 600},
		Layout: VBox{},
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					*dtk.MainPage,
					*hdk.MainPage,
					*ddjb.MainPage,
					*tky.MainPage,
				},
			},
		},
	}).Run(); err != nil {
		fmt.Println(err)
	}
}
