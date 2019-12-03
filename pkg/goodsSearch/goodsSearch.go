package goodsSearch

import (
	"bufio"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wppzxc/taobao_links/pkg/goodsSearch/app"
	"io"
	"os"
	"strings"
	"unicode"
)

type GoodsSearch struct {
	ParentWindow *walk.MainWindow
	MainPage *TabPage
	GetBtn *walk.PushButton
	FileName *walk.LineEdit
	InputFile string
}
var goodsSearch = &GoodsSearch{}
func GetGoodsSearchPage() *GoodsSearch {
	goodsSearch.MainPage = &TabPage{
		Title: "TB商品销量搜索",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: goodsSearch,
			AutoSubmit: true,
		},
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 200},
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "选择文件",
						OnClicked: goodsSearch.ChooseFile,
					},
					LineEdit{
						ReadOnly: true,
						AssignTo: &goodsSearch.FileName,
					},
				},
			},
			Composite{
				MaxSize: Size{0, 200},
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "生成并导出",
						OnClicked: goodsSearch.SearchGoods,
						AssignTo: &goodsSearch.GetBtn,
					},
				},
			},
		},
	}
	return goodsSearch
}

func (g *GoodsSearch) ChooseFile() {
	dlg := new(walk.FileDialog)
	dlg.FilePath = g.InputFile
	dlg.Filter = "文本文档 (*.txt)"
	dlg.Title = "选择文件"
	if ok, err := dlg.ShowOpen(g.ParentWindow); err != nil {
		fmt.Printf("Error in open file : %s\n", err)
	} else if !ok {
		return
	}
	
	g.InputFile = dlg.FilePath
	fmt.Println(g.InputFile)
	g.FileName.SetText(g.InputFile)
	//te := new(walk.Dialog)
	//te.SetTitle("提示")
	//te.SetName("Name")
	//te.Show()
}

func (g *GoodsSearch) SearchGoods() {
	g.SetUIEnable(false)
	file, _ := os.Open(g.InputFile)
	defer file.Close()
	rd := bufio.NewReader(file)
	titles := []string{}
	for {
		line, err := rd.ReadString('\n')
		if io.EOF == err {
			fmt.Println("加载文件完毕")
			titles = append(titles, line)
			break
		}
		if len(line) == 0 {
			continue
		}
		titles = append(titles, line)
	}
	finalTitles := []string{}
	for _, t := range titles {
		tt := strings.TrimFunc(t, unicode.IsSpace)
		finalTitles = append(finalTitles, tt)
	}
	go func() {
		if err := app.SearchAndSave(finalTitles); err != nil {
			fmt.Printf("Error in search goods : %s\n", err)
		}
		g.SetUIEnable(true)
	}()
}

func (g *GoodsSearch) SetUIEnable(enable bool) {
	g.GetBtn.SetEnabled(enable)
	g.FileName.SetEnabled(enable)
}