package duoduojinbao

import (
	"fmt"
	"github.com/wpp/taobao_links/pkg/duoduojinbao/javascripts"
	"github.com/wpp/taobao_links/pkg/duoduojinbao/utils"
	"github.com/zserge/webview"
	"strings"
	"time"
)

type Application struct {
	WebApp    webview.WebView
	Logined   bool
	AccessKey string
}

var logined = false
var AK = ""

func StartLogin(d *Duoduojinbao) {
	loginPage := webview.New(webview.Settings{
		Title:                  "Login",
		URL:                    "https://jinbao.pinduoduo.com",
		ExternalInvokeCallback: eventHandler,
	})
	login := Application{
		WebApp: loginPage,
	}
	go func() {
		fmt.Println("reset ak ...")
		time.Sleep(1 * time.Second)
		login.RestAK()
		fmt.Println("start login ...")
		for {
			login.Login()
			if logined {
				d.AK = AK
				login.CloseLoginPage()
				return
			}
		}
	}()
	fmt.Println("start login ...")
	login.WebApp.Run()
	login.WebApp.Exit()
}

func eventHandler(w webview.WebView, data string) {
	strs := strings.Split(data, "|||")
	fmt.Println("event is : ", strs[0])
	switch strs[0] {
	case "cookie":
		cookies := utils.ParseData(data)
		ak, ok := cookies["PDDAccessToken"]
		if ok {
			fmt.Println("-------------------------AK is :", ak)
			AK = ak
			logined = true
			return
		}
	}
}

func (app *Application) RestAK() {
	app.WebApp.Dispatch(func() {
		if err := app.WebApp.Eval(javascripts.ResetAKJS); err != nil {
			fmt.Printf("Error : %s\n",err)
		}
		fmt.Println("初始化AK成功!")
	})
}

func (app *Application) Login() {
	time.Sleep(1 * time.Second)
	app.WebApp.Dispatch(func() {
		if err := app.WebApp.Eval(javascripts.LoginJS); err != nil {
			fmt.Printf("Error : %s\n",err)
		}
	})
}

func (app *Application) CloseLoginPage() {
	app.WebApp.Dispatch(func() {
		if err := app.WebApp.Eval(javascripts.CloseLoginPage); err != nil {
			fmt.Printf("Error : %s\n",err)
		}
		fmt.Println("登录成功!")
	})
}