package goodsSearch

import "testing"

func TestSearchGoods(t *testing.T) {
	goodsSearch := &GoodsSearch{}
	goodsSearch.InputFile = "titles.txt"
	goodsSearch.SearchGoods()
}
