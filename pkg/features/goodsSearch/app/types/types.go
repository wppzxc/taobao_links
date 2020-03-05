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
	SearchStr     string `json:"searchStr"`
	Title         string `json:"title"`
	GoodsNumber   string `json:"goodsNumber"`
	TaobaoNumber  string `json:"shopNumber"`
	MaxSaleNumber string `json:"maxSaleNumber"`
	RebateAmount  string `json:"rebateAmount"`
	HandPrice     string `json:"handPrice"`
	// rate = [RebateAmount/(RebateAmount + HandPrice)] * 2
	Rate         string `json:"rate"`
	GoodsId      string `json:"goodsId"`
	Url          string `json:"url"`
	IsTitleRight string `json:"isTitleRight"`
}
