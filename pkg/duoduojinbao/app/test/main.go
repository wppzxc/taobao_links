package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/wpp/taobao_links/pkg/duoduojinbao/app"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// 将拼多多链接转为在线拼团人数，并保存到本地
func main() {
	var key string
	flag.StringVar(&key, "key", "", "拼多多key")
	flag.Parse()
	if len(key) == 0 {
		fmt.Print("请输入key!")
		os.Exit(0)
	}
	linksFile := "links.txt"
	//linksFile := "D:\\project\\go\\src\\github.com\\wpp\\taobao_links\\pkg\\duoduojinbao\\app\\test\\links.txt"
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
			links = append(links, line)
			break
		} else if err != nil {
			fmt.Printf("读取第%d行数据失败", index)
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
		fmt.Print("读取文件失败！")
		os.Exit(1)
	}
	cels := app.GetMemNumberLinks(finalLinks, key)
	
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
