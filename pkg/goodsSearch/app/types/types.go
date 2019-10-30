package types

type TaoFenBaResult struct {
	GoodsList []GoodsInfo `json:"goodsList"`
}

type GoodsInfo struct {
	Title      string `json:"title"`
	SaleAmount string `json:"saleAmount"`
	ShopName   string `json:"shopName"`
	GoodsId    string `json:"itemId"`
	Refer      string `json:"refer"`
}

type CelData struct {
	SearchStr     string `json:"searchStr"`
	Title         string `json:"title"`
	GoodsNumber   string `json:"goodsNumber"`
	TaobaoNumber  string `json:"shopNumber"`
	MaxSaleNumber string `json:"maxSaleNumber"`
	GoodsId       string `json:"goodsId"`
	Url           string `json:"url"`
	IsTitleRight  string `json:"isTitleRight"`
}
