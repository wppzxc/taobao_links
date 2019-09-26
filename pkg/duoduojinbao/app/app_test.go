package app

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"os"
	"strconv"
	"testing"
	"time"
)

//func TestGetNameMemberByLink(t *testing.T) {
//	url := "https://mobile.yangkeduo.com/goods2.html?goods_id=34563598420"
//	key := ""
//	data, num := GetNameMemberByLink(url, key)
//	fmt.Println(num)
//	fmt.Println(data)
//}

func TestGetNameMemberByLinks(t *testing.T) {
	var key string
	flag.StringVar(&key, "key", "", "拼多多key")
	flag.Parse()
	if len(key) == 0 {
		fmt.Print("请输入key!")
		os.Exit(0)
	}
	linksFile := "./links.txt"
	file, err := os.Open(linksFile)
	if err != nil {
		fmt.Print(err)
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
			fmt.Print("读取数据结束")
			break
		} else if err != nil {
			fmt.Printf("读取第%d行数据失败", index)
			continue
		}
		links = append(links, line)
	}
	if len(links) == 0 {
		fmt.Print("读取文件失败！")
		os.Exit(1)
	}
	cels := GetMemNumberLinks(links, key)
	
	// 保存excel
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
	err = excelFile.Save(filename)
	if err != nil {
		fmt.Printf("保存excel失败：%s", err.Error())
		os.Exit(1)
	}
}
