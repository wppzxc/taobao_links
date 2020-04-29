package app

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/wppzxc/taobao_links/pkg/features/goodsSearch/app/types"
	"github.com/wppzxc/taobao_links/pkg/utils/dataokeapi"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	openapiListSuperGoods = "https://openapi.dataoke.com/api/goods/list-super-goods"
	appKey                = "5e9d2dbadc286"
	appSecret             = "8f3c81484fdf7bd2695ddbbc6a128201"
	apiVersion            = "v1.2.2"
	searchType            = "0"
	totalSalesDes         = "total_sales_des"
	shopTypeTaobao        = 0
	shopTypeTmall         = 1
)

const (
	right    = "是"
	notRight = "否"
)

func SearchAndSave(titles []string) error {
	cels := []types.CelData{}
	get := func(title string) {
		params := url.Values{
			"version":  []string{apiVersion},
			"type":     []string{searchType},
			"pageId":   []string{"1"},
			"pageSize": []string{"100"},
			"sort":     []string{totalSalesDes},
		}
		tmpUrl := fmt.Sprintf("appKey=%s&keyWords=%s&%s", appKey, title, params.Encode())
		sign := dataokeapi.MakeSign(tmpUrl, appSecret)
		params["appKey"] = []string{appKey}
		params["keyWords"] = []string{title}
		params["sign"] = []string{sign}
		reqUrl := fmt.Sprintf("%s?%s", openapiListSuperGoods, params.Encode())
		resp, err := http.Get(reqUrl)
		if err != nil {
			fmt.Println("Error in get list super goods : ", err)
			return
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		result := new(types.DataokeListSuperGoodsResponseData)
		if err := json.Unmarshal(data, result); err != nil {
			fmt.Println("Error in decode data : ", err)
			return
		}
		cel := new(types.CelData)
		cel.SearchStr = title
		if len(result.Data.List) == 100 {
			cel.GoodsNumber = "100"
		} else {
			cel.GoodsNumber = strconv.Itoa(len(result.Data.List))
		}
		m, num := getMaxSaleAndTaobaoNumber(result.Data.List)
		if m != nil {
			cel.TaobaoNumber = num
			cel.Title = m.Title
			cel.ItemLink = m.ItemLink
			cel.DailySales = m.DailySales
			cel.TwoHoursSales = m.TwoHoursSales
			cel.MonthSales = m.MonthSales
			cel.CommissionRate = m.CommissionRate
			cel.ActualPrice = m.ActualPrice
			cel.CouponPrice = m.CouponPrice
			cel.OriginalPrice = m.OriginalPrice
			if cel.Title == cel.SearchStr {
				cel.IsTitleRight = right
			} else {
				cel.IsTitleRight = notRight
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

func getMaxSaleAndTaobaoNumber(goods []types.DataokeGoods) (*types.DataokeGoods, int64) {
	if len(goods) == 0 {
		return nil, 0
	}
	index := 0
	var taobaoNum int64 = 0
	var lastNum int64 = 0
	for i, g := range goods {
		if g.ShopType == shopTypeTaobao {
			taobaoNum++
		}
		if g.MonthSales >= lastNum {
			lastNum = g.MonthSales
			index = i
		}
	}
	return &goods[index], taobaoNum
}

func saveExcel(cels []types.CelData) {
	timestamp := time.Now().Unix()
	filename := "TBGoodsSearch" + strconv.FormatInt(timestamp, 10) + ".xlsx"
	var excelFile *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	excelFile = xlsx.NewFile()
	sheet, _ = excelFile.AddSheet("商品信息")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "搜索的标题"
	cell = row.AddCell()
	cell.Value = "在线商品数量"
	cell = row.AddCell()
	cell.Value = "淘宝店铺数量"
	cell = row.AddCell()
	cell.Value = "商品原价"
	cell = row.AddCell()
	cell.Value = "优惠券金额"
	cell = row.AddCell()
	cell.Value = "券后价"
	cell = row.AddCell()
	cell.Value = "佣金比例"
	cell = row.AddCell()
	cell.Value = "30天销量"
	cell = row.AddCell()
	cell.Value = "当天销量"
	cell = row.AddCell()
	cell.Value = "2小时销量"
	cell = row.AddCell()
	cell.Value = "商品链接"
	cell = row.AddCell()
	cell.Value = "商品标题"
	cell = row.AddCell()
	cell.Value = "是否一致"
	for _, c := range cels {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = c.SearchStr
		cell = row.AddCell()
		cell.Value = c.GoodsNumber
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", c.TaobaoNumber)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%f", c.OriginalPrice)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%f", c.CouponPrice)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%f", c.ActualPrice)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%f", c.CommissionRate)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", c.MonthSales)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", c.DailySales)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", c.TwoHoursSales)
		cell = row.AddCell()
		cell.Value = c.ItemLink
		cell = row.AddCell()
		cell.Value = c.Title
		cell = row.AddCell()
		cell.Value = c.IsTitleRight
	}
	err := excelFile.Save(filename)
	if err != nil {
		fmt.Printf("Error in save excel file : %s\n", err)
	}
}
