package filter

import (
	"encoding/json"
	"fmt"
	"github.com/wppzxc/taobao_links/pkg/features/yituike/types"
	"testing"
)

func TestFilterByQuan(t *testing.T) {
	config := &types.Config{}
	config.Fanli.FilterQuanNum = 100
	data := `[{"id":28812,"activity_status":0,"start_time":1571323200,"stop_time":1571358600,"extend_document":"#拼多多免单来袭请注意看清要求\n###1.领取【1】元券，拍【1瓶装】选项，券后【5.9】元拍下\n###2.实发【1瓶洁厕液不带喷头】\n###3.禁止使用平台券，禁止联系商家咨询\n###4.有任何问题找群主\n###5.券无代表活动结束\n###6.拍完付款后请重新进入活动网址，点订单页面查询是否正常\n###7.【买家必须自己拍照上传产品图，上满6张图或者拍视频上传，最好能追加一下图片，商家不易，请配合】不上图评价的监测到下次不给单了\n###8.【评价优质用户，下次放单优先给单，不合格的拉黑，下次不放，天天有单，合作愉快】","goods_image_url":"http://t00img.yangkeduo.com/goods/images/2019-10-17/98e89b55-356b-42a0-9589-92e877403a8b.jpg","goods_name":"【强效洁厕液】马桶清洁剂 洁厕灵洁厕宝洗厕所洁厕宝尿垢除污","min_group_price":"6.90","coupon_discount":"1.00","refund_amount":"5.90","service_charge":"0.50","coupon_remain_quantity":20,"coupon_total_quantity":20,"order_mode":1,"from":2,"admin_name":"驰天"},{"id":28811,"activity_status":0,"start_time":1571323200,"stop_time":1571409600,"extend_document":"拼多多免单来袭请注意看清要求，不能不要吃，商家不易，不易啊！\n1.领取【1】元券，颜色：白色，组合：简装。券后【4.1】元拍下，返款5.1元。\n2.空包不发货，多返款一元，必须评价，否则不要吃\n3.禁止使用平台券，禁止联系商家咨询\n4.有任何问题找群主\n5.券无代表活动结束\n6.拍完付款后请重新进入活动网址，点订单页面查询是否正常\n7.必须五星好评，文字评价自由发挥，否则查到踢群，拉黑！","goods_image_url":"http://m.wangshikun.wang/2019_10_17_2227156a32f594b699c05896a78166ad98107.jpg?imageMogr2/auto-orient/thumbnail/300x300/format/jpg/blur/1x0/quality/75|imageslim","goods_name":"韩国ins同款秋冬热卖仿兔毛发圈可爱毛绒纯色皮筋马尾发绳扎头发","min_group_price":"5.10","coupon_discount":"1.00","refund_amount":"5.10","service_charge":"0.50","coupon_remain_quantity":100,"coupon_total_quantity":100,"order_mode":1,"from":2,"admin_name":"踏浪"},{"id":28810,"activity_status":0,"start_time":1571323200,"stop_time":1571409600,"extend_document":"拼多多免单来袭请注意看清要求\n1.领取【1】元券，拍【70*100cm】选项，券后【8.8】元拍下\n2.实发【空包】\n3.禁止使用平台券，禁止联系商家咨询\n4.有任何问题找群主\n5.券无代表活动结束\n6.拍完付款后请重新进入活动网址，点订单页面查询是否正常\n7.必须五星好评，文字评价自由发挥，否则查到踢群，拉黑！\n8.【注意】刷点回复率，刷点回复率，刷点回复率","goods_image_url":"http://m.wangshikun.wang/2019_10_17_2224372019_10_17_22529111.jpg?imageMogr2/auto-orient/thumbnail/300x300/format/jpg/blur/1x0/quality/75|imageslim","goods_name":"法兰绒加厚毛毯午休盖毯子学生单人双人法莱绒180克珊瑚绒床单","min_group_price":"9.80","coupon_discount":"1.00","refund_amount":"9.80","service_charge":"0.50","coupon_remain_quantity":200,"coupon_total_quantity":200,"order_mode":1,"from":2,"admin_name":"米多居"}]`
	items := []types.Item{}
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		fmt.Println(err)
	}
	i := FilterByQuan(config, items)
	fmt.Printf("%#v\n", i)
}
