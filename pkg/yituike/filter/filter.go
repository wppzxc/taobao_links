package filter

import "github.com/wpp/taobao_links/pkg/yituike/types"


// filter items by coupon_remain_quantity
func FilterByQuan(config *types.Config, items []types.Item) *types.Item {
	if config.Fanli.FilterQuanNum == 0 {
		return &items[0]
	}
	resultItem := new(types.Item)
	for _, i := range items {
		if i.CouponRemainQuantity >= resultItem.CouponRemainQuantity {
			resultItem = &i
		}
	}
	if resultItem.CouponRemainQuantity >= config.Fanli.FilterQuanNum {
		return resultItem
	}
	return nil
}