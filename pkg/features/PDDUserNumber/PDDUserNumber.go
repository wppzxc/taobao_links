package PDDUserNumber

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wppzxc/taobao_links/pkg/features/PDDUserNumber/app"
)

type PDDUserNumber struct {
	ParentWindow *walk.MainWindow
	PDDKey       string
	MainPage     *TabPage
	FileName     *walk.LineEdit
	InputFile    string
	RunBtn       *walk.PushButton
}

func GetPDDUserNumberPage() *PDDUserNumber {
	pddUserNumber := &PDDUserNumber{}
	pddUserNumber.MainPage = &TabPage{
		Title:  "PDD在线拼团人数查询",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: pddUserNumber,
			AutoSubmit: true,
		},
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 200},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					Label{
						Text: "拼多多token",
					},
					LineEdit{
						Text: Bind("PDDKey"),
					},
				},
			},
			Composite{
				MaxSize: Size{0, 200},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text:      "选择文件",
						OnClicked: pddUserNumber.ChoosFile,
					},
					LineEdit{
						AssignTo: &pddUserNumber.FileName,
						ReadOnly: true,
						Text:     Bind("InputFile"),
					},
				},
			},
			Composite{
				MaxSize: Size{0, 200},
				Layout:  VBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text:      "生成并导出",
						OnClicked: pddUserNumber.Run,
						AssignTo: &pddUserNumber.RunBtn,
					},
				},
			},
		},
	}
	return pddUserNumber
}

func (p *PDDUserNumber) ChoosFile() {
	dlg := new(walk.FileDialog)
	dlg.FilePath = p.InputFile
	dlg.Filter = "文本文件 (*.txt)"
	dlg.Title = "选择文件"
	if ok, err := dlg.ShowOpen(p.ParentWindow); err != nil {
		fmt.Printf("Error in open file : %s\n", err)
	} else if !ok {
		return
	}
	
	p.InputFile = dlg.FilePath
	p.FileName.SetText(p.InputFile)
	fmt.Println(p.InputFile)
}

func (p *PDDUserNumber) Run() {
	if len(p.PDDKey) == 0 {
		fmt.Println("Error in PDDkey : it must be provided")
		return
	}
	if len(p.InputFile) == 0 {
		fmt.Println("Error in InputFile : it must be provided")
		return
	}
	p.SetUIEnable(false)
	go func() {
		app.Start(p.PDDKey, p.InputFile)
		p.SetUIEnable(true)
	}()
	fmt.Println("Run!")
}

func (p *PDDUserNumber) SetUIEnable(enable bool) {
	p.RunBtn.SetEnabled(enable)
	p.FileName.SetEnabled(enable)
}

func (p *PDDUserNumber) GetMainPage() *TabPage {
	return p.MainPage
}