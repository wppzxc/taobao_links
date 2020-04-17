package duoduojinbao

import (
	"fmt"
	"github.com/wppzxc/taobao_links/pkg/features/duoduojinbao/utils"
	"github.com/zserge/lorca"
	"strings"
	"testing"
)

func TestGetMutiPageLinks(t *testing.T) {
	//up := app.UpData{
	//	CategoryId: -1,
	//	PageNumber: 1,
	//	PageSize: 60,
	//	WithCoupon: 0,
	//}
	//duo := Duoduojinbao{
	//	AK:"cd09bebca57d546bb92a194cd95a6b1d",
	//	PDDAK: "K2WQME43FUPO3ELUPDJWIAO7QKTVQBNFPY5R7UMZMCEMWC5QHRHQ11132a5",
	//	//RangeFrom: "10",
	//	//RangeTo: "20",
	//	LeiMu: -1,
	//}
	//text := duo.GetPDDLinksTest()
	//fmt.Println(text)
}

func TestChrome(t *testing.T) {
	ui, _ := lorca.New("https://jinbao.pinduoduo.com", "", 480, 320)
	defer ui.Close()
	<-ui.Done()
}

func TestTranMoneySep(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	str := "123313212123123213"
	fmt.Println(strings.Index(str, "￥"))
	fmt.Println(len(str))
	result := utils.TranMoneySep(str)
	fmt.Println(result)
}