package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
	"github.com/wppzxc/taobao_links/pkg/features/PDDUserNumber/types"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func Start(key string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	links := []string{}
	index := 0
	for {
		index ++
		line, err := rd.ReadString('\n')
		if io.EOF == err {
			fmt.Println("读取数据结束")
			links = append(links, line)
			break
		} else if err != nil {
			fmt.Printf("读取第%d行数据失败\n", index)
			continue
		}
		if len(line) == 0 {
			continue
		}
		links = append(links, line)
	}
	finalLinks := []string{}
	for _, l := range links {
		link := strings.TrimFunc(l, unicode.IsSpace)
		finalLinks = append(finalLinks, link)
	}
	if len(finalLinks) == 0 {
		fmt.Println("读取文件失败！")
		return
	}
	cels := GetMemNumberLinks(finalLinks, key)
	
	// 保存excel
	timestamp := time.Now().Unix()
	exlFilename := "PDDUserNumber" + strconv.FormatInt(timestamp, 10) + ".xlsx"
	var excelFile *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	excelFile = xlsx.NewFile()
	sheet, _ = excelFile.AddSheet("商品信息")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "在线拼单人数"
	cell = row.AddCell()
	cell.Value = "商品链接"
	cell = row.AddCell()
	cell.Value = "商品名称"
	for _, r := range cels {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", r.LocalGroupTotal)
		cell = row.AddCell()
		cell.Value = r.GoodsLink
		cell = row.AddCell()
		cell.Value = r.GoodsName
	}
	err = excelFile.Save(exlFilename)
	if err != nil {
		fmt.Printf("保存excel失败：%s\n", err.Error())
		os.Exit(1)
	}
}

func GetMemNumberLinks(links []string, key string) []types.ExcelData {
	result := []types.ExcelData{}
	for _, l := range links {
		c := types.ExcelData{
			GoodsLink: l,
		}
		cel, err := GetRawDataByCel(c, key)
		if err != nil {
			continue
		}
		result = append(result, types.ExcelData{
			LocalGroupTotal: cel.Store.InitDataObj.Goods.NeighborGroup.NeighborData.LocalGroup.LocalGroupTotal,
			GoodsLink:       l,
			GoodsName:       cel.Store.InitDataObj.Goods.GoodsName,
		})
	}
	return result
}

func GetRawDataByCel(cel types.ExcelData, key string) (*types.RawData, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, cel.GoodsLink, nil)
	cookie := &http.Cookie{Name: "PDDAccessToken", Value: key}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dom, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	data := ""
	dom.Find("script").Each(func(i int, content *goquery.Selection) {
		str := content.Text()
		strs := strings.Split(str, "window.rawData=")
		if len(strs) > 1 {
			data = strs[1]
		}
	})
	data = strings.TrimFunc(data, unicode.IsSpace)
	data = data[:len(data)-1]
	rd := &types.RawData{}
	if err := json.Unmarshal([]byte(data), rd); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return rd, nil
}
