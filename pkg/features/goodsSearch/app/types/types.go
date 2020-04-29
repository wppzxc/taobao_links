package types

type TaoFenBaResult struct {
	GoodsList []GoodsInfo `json:"goodsList"`
}

type GoodsInfo struct {
	Title        string `json:"title"`
	SaleAmount   string `json:"saleAmount"`
	ShopName     string `json:"shopName"`
	GoodsId      string `json:"itemId"`
	Refer        string `json:"refer"`
	CouponAmount string `json:"couponAmount"`
	RebateAmount string `json:"rebateAmount"`
	HandPrice    string `json:"handPrice"`
}

type CelData struct {
	SearchStr    string `json:"searchStr"`
	Title        string `json:"title"`
	GoodsNumber  string `json:"goodsNumber"`
	TaobaoNumber int64  `json:"shopNumber"`
	// 链接
	ItemLink string `json:"itemLink"`
	// 当天销量
	DailySales int64 `json:"dailySales"`
	// 2小时销量
	TwoHoursSales int64 `json:"twoHoursSales"`
	// 30天销量
	MonthSales int64 `json:"monthSales"`
	// 佣金比例
	CommissionRate float32 `json:"commissionRate"`
	// 券后价
	ActualPrice float32 `json:"actualPrice"`
	// 优惠券金额
	CouponPrice float32 `json:"couponPrice"`
	// 商品原价
	OriginalPrice float32 `json:"originalPrice"`
	IsTitleRight  string  `json:"isTitleRight"`
}

type DataokeListSuperGoodsResponseData struct {
	Time int64       `json:"time"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data DataokeData `json:"data"`
}

type DataokeData struct {
	TotalNum int64          `json:"totalNum"`
	PageId   string         `json:"pageId"`
	List     []DataokeGoods `json:"list"`
}

type DataokeGoods struct {
	// 标题
	Title string `json:"title"`
	// 链接
	ItemLink string `json:"itemLink"`
	// 当天销量
	DailySales int64 `json:"dailySales"`
	// 2小时销量
	TwoHoursSales int64 `json:"twoHoursSales"`
	// 30天销量
	MonthSales int64 `json:"monthSales"`
	// 佣金比例
	CommissionRate float32 `json:"commissionRate"`
	// 券后价
	ActualPrice float32 `json:"actualPrice"`
	// 优惠券金额
	CouponPrice float32 `json:"couponPrice"`
	// 商品原价
	OriginalPrice float32 `json:"originalPrice"`
	// 店铺类型 1-天猫，0-淘宝
	ShopType int `json:"shopType"`
}
