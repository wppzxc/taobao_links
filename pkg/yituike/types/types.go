package types

const (
	AuthUrl = "http://v2.yituike.com/index/login/indexlogin"
	PreUrl  = "http://v2.yituike.com/admin/goods/goods_store?limit=100&page=1&state=2"
	ProUrl  = "http://v2.yituike.com/admin/goods/goods_store?limit=100&page=1&state=1"
)

type Config struct {
	Auth          AuthInfo  `json:"auth"`
	ReceiversFile string    `json:"receiversFile"`
	Receivers     []RecInfo `json:"receivers"`
	Fanli         FanliInfo `json:"fanli"`
}

type AuthInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Url      string `json:"url"`
}

type FanliInfo struct {
	RefreshInterval int64      `json:"refreshInterval"`
	SendInterval    int64      `json:"sendInterval"`
	FilterQuanNum   int        `json:"filterQuanNum"`
	Process         Process    `json:"process"`
	Premonitor      Premonitor `json:"premonitor"`
}

type Process struct {
	Url       string `json:"url"`
	Start     bool   `json:"start"`
	MsgPrefix string `json:"msgPrefix"`
}

type Premonitor struct {
	Url       string `json:"url"`
	Start     bool   `json:"start"`
	MsgPrefix string `json:"msgPrefix"`
}

type RecInfo struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
	SkipImage bool   `json:"skipImage"`
}

// {"code":1,"msg":"暂无数据！","count":"","data":""}
type ItemResult struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Count int    `json:"count"`
	Data  []Item `json:"data,omitempty"`
}

type Item struct {
	Id                   int    `json:"id"`
	ActivityStatus       int    `json:"activity_status"`
	StartTime            int64  `json:"start_time"`
	StopTime             int64  `json:"stop_time"`
	ExtendDocument       string `json:"extend_document"`
	GoodsImageUrl        string `json:"goods_image_url"`
	GoodsName            string `json:"goods_name"`
	MinGroupPrice        string `json:"min_group_price"`
	CouponDiscount       string `json:"coupon_discount"`
	RefundAmount         string `json:"refund_amount"`
	ServiceCharge        string `json:"service_charge"`
	CouponRemainQuantity int    `json:"coupon_remain_quantity"`
	CouponTotalQuantity  int    `json:"coupon_total_quantity"`
	OrderMode            int    `json:"order_mode"`
	From                 int    `json:"from"`
	AdminName            string `json:"admin_name"`
}

type TokenResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data []TokenData `json:"data"`
}

type TokenData struct {
	Token string `json:"token"`
}
