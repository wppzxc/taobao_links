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
	CategoryId int     `json:"categoryId"`
	PageNumber int     `json:"pageNumber"`
	PageSize   int     `json:"pageSize"`
	WithCoupon int     `json:"withCoupon"`
	SortType   int     `json:"sortType"`
	RangeList  []Range `json:"rangeList"`
}

type Range struct {
	RangeFrom int64 `json:"rangeFrom"`
	RangeId   int   `json:"rangeId"`
	RangeTo   int64 `json:"rangeTo"`
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
	MallName string `json:"mallName"`
	// 类目
	CategoryName string `json:"categoryName"`
	// 折扣，需要除100
	CouponDiscount int64 `json:"couponDiscount"`
	// 优惠券剩余
	CouponRemainQuantity int64 `json:"couponRemainQuantity"`
	// 原价，需要除1000
	MinGroupPrice int64 `json:"minGroupPrice"`
	// 销量
	SalesTip string `json:"salesTip"`
	// 佣金比率
	PromotionRate int64  `json:"promotionRate"`
	GoodsId       int64  `json:"goodsId"`
	GoodsName     string `json:"goodsName"`
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
	MallName        string `json:"mallName"`
	// 类目
	CategoryName string `json:"categoryName"`
	// 折扣，需要除100
	CouponDiscount int64 `json:"couponDiscount"`
	// 优惠券剩余
	CouponRemainQuantity int64 `json:"couponRemainQuantity"`
	// 券后价（等于原价减去折扣）
	MinGroupPrice string `json:"minGroupPrice"`
	// 销量
	SalesTip string `json:"salesTip"`
	// 佣金比率
	PromotionRate string `json:"promotionRate"`
	// 佣金
	Promotion string `json:"PromotionRate"`
}

func GetOnePageLinks(upData UpData, key string) []ExcelData {
	reqBody, _ := json.Marshal(upData)
	fmt.Println(reqBody)
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
	cels := []ExcelData{}
	for _, i := range d.Result.GoodsList {
		cels = append(cels, ExcelData{
			GoodsLink:            fmt.Sprintf(getPDDLinksPrefix, i.GoodsId),
			MallName:             i.MallName,
			CategoryName:         i.CategoryName,
			CouponDiscount:       i.CouponDiscount / 1000,
			CouponRemainQuantity: i.CouponRemainQuantity,
			MinGroupPrice:        fmt.Sprintf("%.2f", float32(i.MinGroupPrice)/1000-float32(i.CouponDiscount)/1000),
			SalesTip:             i.SalesTip,
			PromotionRate:        fmt.Sprintf("%d%%", i.PromotionRate/10),
			Promotion:            fmt.Sprintf("%.2f", (float32(i.MinGroupPrice)/1000-float32(i.CouponDiscount)/1000)*(float32(i.PromotionRate)/1000)),
		})
	}
	return cels
}

func GetTextFromLinks(cels []ExcelData) string {
	text := ""
	for _, c := range cels {
		text = text + fmt.Sprintf("%s----%d", c.GoodsLink, c.LocalGroupTotal) + "\n"
	}
	return text
}

func GetMemNumberCels(cels []ExcelData, key string) []ExcelData {
	result := []ExcelData{}
	for _, c := range cels {
		rawData, err := GetRawDataByCel(c, key)
		if err != nil {
			continue
		}
		c.LocalGroupTotal = rawData.Store.InitDataObj.Goods.NeighborGroup.NeighborData.LocalGroup.LocalGroupTotal
		c.GoodsName = rawData.Store.InitDataObj.Goods.GoodsName
		result = append(result, c)
	}
	return result
}

func GetMemNumberLinks(links []string, key string) []ExcelData {
	result := []ExcelData{}
	for _, l := range links {
		c := ExcelData{
			GoodsLink: l,
		}
		cel, err := GetRawDataByCel(c, key)
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

func GetRawDataByCel(cel ExcelData, key string) (*RawData, error) {
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
	rd := &RawData{}
	if err := json.Unmarshal([]byte(data), rd); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return rd, nil
}
