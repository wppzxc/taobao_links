package taokeyi

import (
	. "github.com/lxn/walk/declarative"
)

func GetTaokeyiPage() TabPage {
	return TabPage{
		Title:  "淘客易",
		Layout: VBox{},
		Children: []Widget{},
	}
}
