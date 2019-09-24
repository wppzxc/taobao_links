package app

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"
)

const (
	getDuoDataUrl       = "https://jinbao.pinduoduo.com/network/api/common/goodsList"
	getDuoDataCookieKey = "DDJB_PASS_ID"
	getPDDLinksPrefix   = "https://mobile.yangkeduo.com/goods2.html?goods_id=%d"
)

type UpData struct {
	CategoryId int `json:"categoryId"`
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	WithCoupon int `json:"withCoupon"`
	SortType   int `json:"sortType"`
}

type DuoDuoData struct {
	Success   bool   `json:"success"`
	ErrorCode int64  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Result    Result `json:"result"`
}

type Result struct {
	Total     int64  `json:"total"`
	GoodsList []Item `json:"goodsList"`
}

type Item struct {
	MallName  string `json:"mallName"`
	GoodsId   int64  `json:"goodsId"`
	GoodsName string `json:"goodsName"`
}

type RawData struct {
	Store Store `json:"store"`
}

type Store struct {
	InitDataObj InitDataObj `json:"initDataObj"`
}

type InitDataObj struct {
	Goods Goods `json:"goods"`
}

type Goods struct {
	GoodsName     string        `json:"goodsName"`
	NeighborGroup NeighborGroup `json:"neighborGroup"`
}

type NeighborGroup struct {
	NeighborStatus int          `json:"neighbor_status"`
	NeighborData   NeighborData `json:"neighbor_data"`
}

type NeighborData struct {
	LocalGroup LocalGroup `json:"local_group"`
}

type LocalGroup struct {
	LocalGroupTotal int    `json:"local_group_total"`
	LocalGroupDesc  string `json:"local_group_desc"`
}

type ExcelData struct {
	LocalGroupTotal int
	GoodsLink       string
	GoodsName       string
}

func GetOnePageLinks(upData UpData, key string) []string {
	reqBody, _ := json.Marshal(upData)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, getDuoDataUrl, strings.NewReader(string(reqBody)))
	cookie := &http.Cookie{Name: getDuoDataCookieKey, Value: key, HttpOnly: false}
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	d := &DuoDuoData{}
	if err := json.Unmarshal(data, d); err != nil {
		fmt.Println(err)
		return nil
	}
	links := []string{}
	for _, i := range d.Result.GoodsList {
		links = append(links, fmt.Sprintf(getPDDLinksPrefix, i.GoodsId))
	}
	return links
}

func GetTextFromLinks(cels []ExcelData) string {
	text := ""
	for _, c := range cels {
		text = text + fmt.Sprintf("%s----%d", c.GoodsLink, c.LocalGroupTotal) + "\n"
	}
	return text
}

func GetMemNumberLinks(links []string, key string) []ExcelData {
	result := []ExcelData{}
	for _, l := range links {
		cel, err := GetNameMemberByLink(l, key)
		fmt.Printf("%#V", cel)
		if err != nil {
			continue
		}
		result = append(result, ExcelData{
			LocalGroupTotal: cel.Store.InitDataObj.Goods.NeighborGroup.NeighborData.LocalGroup.LocalGroupTotal,
			GoodsLink:       l,
			GoodsName:       cel.Store.InitDataObj.Goods.GoodsName,
		})
	}
	return result
}

func GetNameMemberByLink(link string, key string) (*RawData, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, link, nil)
	cookie := &http.Cookie{Name: "PDDAccessToken", Value: key}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	dom, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Print(err)
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
	cel := &RawData{}
	if err := json.Unmarshal([]byte(data), cel); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return cel, nil
}
