package runtime

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	uuid "github.com/satori/go.uuid"
	"github.com/wppzxc/taobao_links/pkg/version"
	"gopkg.in/gomail.v2"
	"runtime"
	"strconv"
	"strings"
)

const (
	htmlSpec = "<br/>"
	htmlSpace = "&nbsp;"
	html4Space = htmlSpace + htmlSpace + htmlSpace + htmlSpace
)

func HandleCrash(mw *walk.MainWindow) {
	if r := recover(); r != nil {
		errId := uuid.NewV4()
		errStr, stacktrace := getPanicInfo(r)
		errMsg := fmt.Sprintf("程序运行异常！\n\n%s\n\n错误报告ID：%s\n是否发送错误报告？", errStr, errId.String())
		result := walk.MsgBox(mw, "错误", errMsg, walk.MsgBoxOKCancel)
		switch result {
		case win.IDOK:
			msg := formatStacktraceToMail(errStr, stacktrace)
			sendMail(errId.String(), msg)
		case win.IDCANCEL:
			fmt.Printf("不发送")
		default:
			fmt.Println("关闭")
		}
	}
}

func getPanicInfo(r interface{}) (string, []byte) {
	const size = 64 << 10
	stacktrace := make([]byte, size)
	stacktrace = stacktrace[:runtime.Stack(stacktrace, false)]
	if _, ok := r.(string); ok {
		return fmt.Sprintf("Observed a panic: %s", r), stacktrace
	} else {
		return fmt.Sprintf("Observed a panic: %#v (%v)", r, r), stacktrace
	}
}

func sendMail(errId string, msg string) {
	mailConn := map[string]string {
		"user": "bachinanfei@163.com",
		"pass": "Aa123456",
		"host": "smtp.163.com",
		"port": "25",
	}
	port, _ := strconv.Atoi(mailConn["port"])

	m := gomail.NewMessage()
	m.SetHeader("From","taobao_links" + "<" + mailConn["user"] + ">")
	m.SetHeader("To", "bachinanfei@qq.com")
	m.SetHeader("Subject", fmt.Sprintf("taobao_links程序崩溃: %s", errId))
	m.SetBody("text/html", msg)

	gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"]).DialAndSend(m)
}

func formatStacktraceToMail(errStr string, stacktrace []byte) string {
	str := string(stacktrace)
	str = strings.ReplaceAll(str, "\t", html4Space)
	strs := strings.Split(str, "\n")
	var result string
	result = result + "version: " + version.Get().String() + htmlSpec
	result = result + errStr + htmlSpec
	for _, s := range strs {
		result = result + s + htmlSpec
	}
	return result
}