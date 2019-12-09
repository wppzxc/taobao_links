package duoduojinbao

import (
	"fmt"
	"github.com/zserge/lorca"
	"testing"
)

func TestGetMutiPageLinks(t *testing.T) {
	//up := app.UpData{
	//	CategoryId: -1,
	//	PageNumber: 1,
	//	PageSize: 60,
	//	WithCoupon: 0,
	//}
	duo := Duoduojinbao{
		AK:"cd09bebca57d546bb92a194cd95a6b1d",
		PDDAK: "K2WQME43FUPO3ELUPDJWIAO7QKTVQBNFPY5R7UMZMCEMWC5QHRHQ11132a5",
		//RangeFrom: "10",
		//RangeTo: "20",
		LeiMu: -1,
	}
	text := duo.GetPDDLinksTest()
	fmt.Println(text)
}

func TestChrome(t *testing.T) {
	ui, _ := lorca.New("https://jinbao.pinduoduo.com", "", 480, 320)
	defer ui.Close()
	<-ui.Done()
}