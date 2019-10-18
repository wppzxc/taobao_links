package duoduojinbao

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tealeg/xlsx"
	"github.com/wpp/taobao_links/pkg/duoduojinbao/app"
	"strconv"
	"time"
)

const (
	JingXuan  = -1
	TuiJian   = -2
	QingCang  = -11
	BaiHuo    = 15
	MuYing    = 4
	ShiPin    = 1
	NvZhuang  = 14
	DianQi    = 18
	XieBao    = 1281
	NeiYi     = 1282
	MeiZhuang = 16
	NanXie    = 743
	ShuiGuo   = 13
	JiaFang   = 818
	WenJu     = 2478
	YunDong   = 1451
	XuNi      = 590
	QiChe     = 2048
	JiaZhuang = 1917
	JiaJu     = 2974
	YiYao     = 3279
)

/*
{
	"categoryId":-1, 精选、推荐...
	"pageNumber": 1, 分页
	"pageSize": 60, 分页
	"withCoupon": 1, 是否有优惠券，有1，没有0
	"sortType":6 销量排序，5降序，6升序
}
*/
type Duoduojinbao struct {
	MainPage    *TabPage
	GetBtn      *walk.PushButton
	DownLoadBtn *walk.PushButton
	Links       *walk.TextEdit
	StartPage   string
	EndPage     string
	LeiMu       int
	Quan        bool
	XiaoLiang   int
	AK          string
	PDDAK       string
	Excel       []app.ExcelData
	RangeFrom   string
	RangeTo     string
}

type LeiMu struct {
	Id   int
	Name string
}

func GetLeiMus() []LeiMu {
	return []LeiMu{
		{Name: "精选", Id: -1},
		{Name: "推荐", Id: -2},
		{Name: "清仓", Id: -11},
		{Name: "百货", Id: 15},
		{Name: "母婴", Id: 4},
		{Name: "饰品", Id: 1},
		{Name: "女装", Id: 14},
		{Name: "电器", Id: 18},
		{Name: "鞋包", Id: 1281},
		{Name: "内衣", Id: 1282},
		{Name: "美妆", Id: 16},
		{Name: "男鞋", Id: 743},
		{Name: "水果", Id: 13},
		{Name: "家纺", Id: 818},
		{Name: "文具", Id: 2478},
		{Name: "运动", Id: 1451},
		{Name: "虚拟", Id: 590},
		{Name: "汽车", Id: 2048},
		{Name: "家装", Id: 1917},
		{Name: "家具", Id: 2974},
		{Name: "医药", Id: 3279},
	}
}

func GetDuoduojinbaoPage() *Duoduojinbao {
	duoduojinbao := &Duoduojinbao{
		Links: &walk.TextEdit{},
	}
	duoduojinbao.MainPage = &TabPage{
		Title:  "多多进宝",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: duoduojinbao,
			AutoSubmit: true,
		},
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 200},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					Label{
						Text: "多多进宝 token ：",
					},
					LineEdit{
						Text: Bind("AK"),
					},
				},
			},
			Composite{
				MaxSize: Size{0, 200},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					Label{
						Text: "拼多多 token ：",
					},
					LineEdit{
						Text: Bind("PDDAK"),
					},
				},
			},
			Composite{
				MaxSize: Size{0, 200},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					Label{
						Text: "类目",
					},
					ComboBox{
						Value:         Bind("LeiMu"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         GetLeiMus(),
					},
					Label{
						Text: "销量排序",
					},
					ComboBox{
						Value:         Bind("XiaoLiang"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         []LeiMu{{Name: "降序", Id: 6}, {Name: "升序", Id: 5}},
					},
					CheckBox{
						Text:           "      优惠券",
						Checked:        Bind("Quan"),
						TextOnLeftSide: true,
					},
					TextLabel{
						Text: "价格（元）",
					},
					LineEdit{
						MinSize: Size{60, 0},
						Text:    Bind("RangeFrom"),
					},
					TextLabel{
						Text: "~",
					},
					LineEdit{
						MinSize: Size{60, 0},
						Text:    Bind("RangeTo"),
					},
				},
			},
			Composite{
				Layout:  HBox{},
				MaxSize: Size{0, 200},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text:      "拉取",
						AssignTo:  &duoduojinbao.GetBtn,
						OnClicked: duoduojinbao.GetPDDLinks,
					},
					PushButton{
						Text:      "导出",
						AssignTo:  &duoduojinbao.DownLoadBtn,
						OnClicked: duoduojinbao.DownLoadXLSX,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "拼多多链接：",
					},
					TextEdit{
						ReadOnly: true,
						AssignTo: &duoduojinbao.Links,
						VScroll:  true,
					},
				},
			},
			Composite{
				Layout: HBox{Margins: Margins{}},
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "起始页：",
					},
					LineEdit{
						MaxSize: Size{20, 0},
						Text:    Bind("StartPage"),
					},
					TextLabel{
						Text: "结束页：",
					},
					LineEdit{
						MaxSize: Size{20, 0},
						Text:    Bind("EndPage"),
					},
				},
			},
		},
	}
	return duoduojinbao
}

func (d *Duoduojinbao) GetPDDLinks() {
	d.SetUIEnable(false)
	upData := app.UpData{
		CategoryId: d.LeiMu,
		PageSize:   60,
		WithCoupon: 0,
		SortType:   d.XiaoLiang,
		RangeList:  []app.Range{},
	}
	rngF, _ := strconv.ParseInt(d.RangeFrom, 10, 64)
	rngT, _ := strconv.ParseInt(d.RangeTo, 10, 64)
	rng := app.Range{
		RangeFrom: rngF * 100,
		RangeId:   1,
		RangeTo:   rngT * 100,
	}
	if rng.RangeFrom > 0 && rng.RangeTo > 0 {
		upData.RangeList = append(upData.RangeList, rng)
	}
	if d.Quan {
		upData.WithCoupon = 1
	}
	startPage, _ := strconv.Atoi(d.StartPage)
	endPage, _ := strconv.Atoi(d.EndPage)
	if startPage <= 0 {
		startPage = 1
	}
	if startPage > endPage {
		endPage = startPage
	}
	fmt.Printf("updata is %#v", upData)
	go func() {
		text := d.GetMutiPageLinks(upData, startPage, endPage)
		d.UpdateLinks(text)
	}()
}

func (d *Duoduojinbao) SetUIEnable(enable bool) {
	d.GetBtn.SetEnabled(enable)
	d.Links.SetEnabled(enable)
	d.DownLoadBtn.SetEnabled(enable)
}

func (d *Duoduojinbao) GetMutiPageLinks(upData app.UpData, startPage int, endPage int) string {
	cels := []app.ExcelData{}
	for i := startPage; i <= endPage; i++ {
		upData.PageNumber = i
		tmpCels := app.GetOnePageLinks(upData, d.AK)
		cels = append(cels, tmpCels...)
	}
	cels = app.GetMemNumberCels(cels, d.PDDAK)
	d.Excel = cels
	text := app.GetTextFromLinks(cels)
	return text
}

func (d *Duoduojinbao) UpdateLinks(text string) {
	defer d.SetUIEnable(true)
	d.Links.SetText(text)
}

func (d *Duoduojinbao) DownLoadXLSX() {
	d.SetUIEnable(false)
	go func() {
		d.SaveXLSX()
	}()
}

func (d *Duoduojinbao) SaveXLSX() {
	defer d.SetUIEnable(true)
	if len(d.Excel) == 0 {
		return
	}
	timestamp := time.Now().Unix()
	filename := "DuoDuojinbao" + strconv.FormatInt(timestamp, 10) + ".xlsx"
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("商品信息")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "商品名称"
	cell = row.AddCell()
	cell.Value = "类目"
	cell = row.AddCell()
	cell.Value = "在线拼单人数"
	cell = row.AddCell()
	cell.Value = "折扣"
	cell = row.AddCell()
	cell.Value = "券后价格"
	cell = row.AddCell()
	cell.Value = "佣金"
	cell = row.AddCell()
	cell.Value = "佣金率"
	cell = row.AddCell()
	cell.Value = "销量"
	cell = row.AddCell()
	cell.Value = "优惠券剩余"
	cell = row.AddCell()
	cell.Value = "商品链接"
	cell = row.AddCell()
	cell.Value = "店铺名称"
	for _, r := range d.Excel {
		row = sheet.AddRow()
		// 商品名称
		cell = row.AddCell()
		cell.Value = r.GoodsName
		// 类目
		cell = row.AddCell()
		cell.Value = r.CategoryName
		// 在线拼单人数
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", r.LocalGroupTotal)
		// 折扣
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", r.CouponDiscount)
		// 券后价格
		cell = row.AddCell()
		cell.Value = r.MinGroupPrice
		// 佣金
		cell = row.AddCell()
		cell.Value = r.Promotion
		// 佣金率
		cell = row.AddCell()
		cell.Value = r.PromotionRate
		// 销量
		cell = row.AddCell()
		cell.Value = r.SalesTip
		// 优惠券剩余
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", r.CouponRemainQuantity)
		// 商品链接
		cell = row.AddCell()
		cell.Value = r.GoodsLink
		// 店铺名称
		cell = row.AddCell()
		cell.Value = r.MallName
	}
	err = file.Save(filename)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
