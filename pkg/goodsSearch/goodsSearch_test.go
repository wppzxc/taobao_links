package goodsSearch

import "testing"

func TestSearchGoods(t *testing.T) {
	goodsSearch := &GoodsSearch{}
	goodsSearch.InputFile = "100.txt"
	goodsSearch.SearchGoods()
}
