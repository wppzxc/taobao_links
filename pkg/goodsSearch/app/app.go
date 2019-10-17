package app

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/wpp/taobao_links/pkg/goodsSearch/app/types"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	taofenbaUrl = "https://www.taofen8.com/search/s"
	taobao = "taobao"
	tmall = "tmall"
	taobaoPrefix = "https://item.taobao.com/item.htm?id="
	tmallPrefix = "https://detail.tmall.com/item.htm?id="
)

func SearchAndSave(titles []string) error {
	cels := []types.CelData{}
	get := func(title string) {
		params := url.Values{}
		params.Set("q", title)
		sUrl, _ := url.Parse(taofenbaUrl)
		sUrl.RawQuery = params.Encode()
		resp, _ := http.Post(sUrl.String(), "application/json;charset=UTF-8", nil)
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		result := new(types.TaoFenBaResult)
		json.Unmarshal(data, result)
		cel := new(types.CelData)
		cel.Title = title
		if len(result.GoodsList) == 100 {
			cel.GoodsNumber = "100"
		} else {
			cel.GoodsNumber = strconv.Itoa(len(result.GoodsList))
		}
		m, num := getMaxSaleAndTaobaoNumber(result.GoodsList)
		if m != nil {
			cel.MaxSaleNumber = m.SaleAmount
			cel.GoodsId = m.GoodsId
			cel.TaobaoNumber = num
			if m.Refer == taobao {
				cel.Url = fmt.Sprint(taobaoPrefix, m.GoodsId)
			} else if m.Refer == tmall {
				cel.Url = fmt.Sprint(tmallPrefix, m.GoodsId)
			}
		}
		cels = append(cels, *cel)
	}
	for _, t := range titles {
		if len(t) == 0 {
			continue
		}
		get(t)
	}
	saveExcel(cels)
	fmt.Println(cels)
	return nil
}

func getMaxSaleAndTaobaoNumber(goods []types.GoodsInfo) (*types.GoodsInfo, string) {
	if len(goods) == 0 {
		return nil, "0"
	}
	index := 0
	taobaoNum := 0
	var lastNum int64 = 0
	for i, g := range goods {
		if g.Refer == taobao {
			taobaoNum ++
		}
		n64, _ := strconv.ParseInt(g.SaleAmount, 10, 64)
		if n64 >= lastNum {
			lastNum = n64
			index = i
		}
	}
	return &goods[index], fmt.Sprintf("%d", taobaoNum)
}

func saveExcel(cels []types.CelData) {
	timestamp := time.Now().Unix()
	filename := strconv.FormatInt(timestamp, 10) + ".xlsx"
	var excelFile *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	excelFile = xlsx.NewFile()
	sheet, _ = excelFile.AddSheet("商品信息")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "商品标题"
	cell = row.AddCell()
	cell.Value = "在线商品数量"
	cell = row.AddCell()
	cell.Value = "淘宝店铺数量"
	cell = row.AddCell()
	cell.Value = "最高销量"
	cell = row.AddCell()
	cell.Value = "商品ID"
	cell = row.AddCell()
	cell.Value = "商品链接"
	for _, c := range cels {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = c.Title
		cell = row.AddCell()
		cell.Value = c.GoodsNumber
		cell = row.AddCell()
		cell.Value = c.TaobaoNumber
		cell = row.AddCell()
		cell.Value = c.MaxSaleNumber
		cell = row.AddCell()
		cell.Value = c.GoodsId
		cell = row.AddCell()
		cell.Value = c.Url
	}
	err := excelFile.Save(filename)
	if err != nil {
		fmt.Printf("Error in save excel file : %s", err)
	}
}
