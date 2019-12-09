package types

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
