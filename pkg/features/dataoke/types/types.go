package types

type DataokeQuanApiResponseData struct {
	Time int64               `json:"time"`
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data DataokeQuanDataList `json:"data"`
}

type DataokeQuanDataList struct {
	List []DataokeData `json:"list"`
}

type DataokeTopApiResponseData struct {
	Time int64         `json:"time"`
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []DataokeData `json:"data"`
}

type DataokeGoodsDetailResponseData struct {
	Time int64       `json:"time"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data DataokeData `json:"data"`
}

type DataokeData struct {
	Id       int64  `json:"id"`
	GoodsId  string `json:"goodsId"`
	Title    string `json:"title"`
	ItemLink string `json:"itemLink"`
}
