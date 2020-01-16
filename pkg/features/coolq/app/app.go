package app

import (
	"encoding/json"
	"fmt"
	"github.com/wppzxc/taobao_links/pkg/features/coolq/types"
	"golang.org/x/net/websocket"
	"net"
	"os/exec"
	"strconv"
	"time"
)

func ConnectWebSocket(wsUrl string) (*websocket.Conn, error) {
	return websocket.Dial(wsUrl, "", "http://api.wpp.pro/")
}

func Start(wsUrl string, groups []string, users []string, interval int, tklTitle string, stopCh chan struct{}) error {
	// start websocket client goroutine
	msgs := make(chan types.Message, 4)
	ws, err := ConnectWebSocket(wsUrl)
	if err != nil {
		return fmt.Errorf("链接 websocket 失败 ： %s，请检查url是否正确！", err)
	}
	go func(ws *websocket.Conn, msgs chan types.Message, stopCh chan struct{}) {
		defer ws.Close()
		for {
			select {
			case <-stopCh:
				fmt.Println("websocket client stoped!")
				return
			default:
				var data []byte
				err := websocket.Message.Receive(ws, &data)
				if err != nil {
					fmt.Println("Error in get message : ", err)
					continue
				}
				msg := &types.Message{}
				if err := json.Unmarshal(data, msg); err != nil {
					fmt.Printf("Error in resolve receive data : %s", err)
				}
				fmt.Printf("get message is : %#v \n", msg)
				if ok := checkSendGroups(msg, groups); ok {
					msgs <- *msg
				}
			}
		}
	}(ws, msgs, stopCh)

	// start sender goroutine
	go func(msgs chan types.Message, users []string, interval int, tklTitle string, stopCh chan struct{}) {
		for {
			select {
			case msg := <-msgs:
				if err := CoolQMessageSend(msg, users, interval, tklTitle); err != nil {
					fmt.Println("Error in send meg : ", err)
				} else {
					fmt.Printf("sender message : %#v \n", msg)
				}
			case <-stopCh:
				fmt.Println("sender stoped!")
				return
			}
		}
	}(msgs, users, interval, tklTitle, stopCh)
	return nil
}

func checkSendGroups(msg *types.Message, groups []string) bool {
	groupId := strconv.FormatInt(msg.GroupId, 10)
	for _, g := range groups {
		if groupId == g {
			return true
		}
	}
	return false
}

func StartLocalCoolQ() error {
	cmd := "./coolq/CQA.exe"
	if err := exec.Command(cmd).Start(); err != nil {
		return fmt.Errorf("Error in start coolq : %s \n", err)
	}
	return nil
}

func CheckCoolqLogined() bool {
	_, err := net.DialTimeout("tcp", "0.0.0.0:6700", 3 * time.Second)
	if err != nil {
		return false
	}
	return true
}
