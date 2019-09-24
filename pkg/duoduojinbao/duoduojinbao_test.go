package duoduojinbao

import (
	"fmt"
	"github.com/wpp/taobao_links/pkg/duoduojinbao/app"
	"github.com/zserge/lorca"
	"testing"
)

func TestGetMutiPageLinks(t *testing.T) {
	up := app.UpData{
		CategoryId: -1,
		PageNumber: 1,
		PageSize: 60,
		WithCoupon: 0,
	}
	duo := Duoduojinbao{
		AK:"ee60c4580b7fa330f42ec18dcf711d2f",
		PDDAK: "7JRUYPDYSMRNH24TX4N73UQMKQFUA533LO5IARUSQD7A65NQ76HQ11132a5",
	}
	text := duo.GetMutiPageLinks(up, 1, 1)
	fmt.Println(text)
}

func TestChrome(t *testing.T) {
	ui, _ := lorca.New("https://jinbao.pinduoduo.com", "", 480, 320)
	defer ui.Close()
	<-ui.Done()
}