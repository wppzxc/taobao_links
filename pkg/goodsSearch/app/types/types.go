package types

type TaoFenBaResult struct {
	GoodsList []GoodsInfo `json:"goodsList"`
}

type GoodsInfo struct {
	Title      string `json:"title"`
	SaleAmount string `json:"saleAmount"`
	ShopName   string `json:"shopName"`
	GoodsId    string `json:"itemId"`
}

type CelData struct {
	Title         string `json:"title"`
	Number        string `json:"number"`
	MaxSaleNumber string `json:"maxSaleNumber"`
	GoodsId       string `json:"goodsId"`
	Url           string `json:"url"`
}
