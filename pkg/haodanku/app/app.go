package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	linkPrefix = "https://detail.tmall.com/item.htm?id="
)

type Result struct {
	Status           int        `json:"status"`
	Message          string     `json:"message"`
	NumPage          int        `json:"num_page"`
	NowTime          uint       `json:"now_time"`
	ItemInfo         []Item     `json:"item_info"`
	ActivityFiltrate []Actflter `json:"activity_filtrate"`
}

type Item struct {
	Id             string `json:"id"`
	ItemId         string `json:"itemid"`
	SellerId       string `json:"seller_id"`
	ItemTitle      string `json:"itemtitle"`
	ItemshortTitle string `json:"itemshorttitle"`
	ItemSale       string `json:"itemsale"`
	ItemSale2      string `json:"itemsale2"`
	ItemPrice      string `json:"itemprice"`
	ItemDesc       string `json:"itemdesc"`
	ItemPicCopy    string `json:"itempic_copy"`
	ShopType       string `json:"shoptype"`
	IsForeShow     string `json:"is_foreshow"`
	IsBrand        string `json:"is_brand"`
	VideoId        string `json:"videoid"`
	ActivityId     string `json:"activityid"`
	ActivityType   string `json:"activity_type"`
	Fqcat          string `json:"fqcat"`
	ShopId         string `json:"shopid"`
	Userid         string `json:"userid"`
	ShopName       string `json:"shopname"`
	Tkrates        string `json:"tkrates"`
	Tkmoney        string `json:"tkmoney"`
	Couponurl      string `json:"couponurl"`
	CouponMoney    string `json:"couponmoney"`
	CouponSurplus  string `json:"couponsurplus"`
	CouponReceive  string `json:"couponreceive"`
	CouponNum      string `json:"couponnum"`
	StartTime      string `json:"starttime"`
	ItemEndPrice   string `json:"itemendprice"`
	GeneralIndex   string `json:"general_index"`
	GradeAvg       string `json:"grade_avg"`
	RateTotal      string `json:"rate_total"`
	TodaySale      string `json:"todaysale"`
	Stock          string `json:"stock"`
	ItemPic        string `json:"itempic"`
	Discount       string `json:"discount"`
	SellerName     string `json:"seller_name"`
	SellerQQ     string `json:"seller_qq"`
	ReportStatus   string `json:"report_status"`
	DownType       string `json:"down_type"`
	ActivityPlan   string `json:"activity_plan"`
}

type Actflter struct {
	ActivityId   string `json:"activity_id"`
	ActivityImg  string `json:"activity_img"`
	ActivityName string `json:"activity_name"`
}

func GetLinks(url string) ([]string, error) {
	var items []string
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	result := Result{}
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, i := range result.ItemInfo {
		link := linkPrefix + i.ItemId
		items = append(items, link)
	}
	return items, nil
}
