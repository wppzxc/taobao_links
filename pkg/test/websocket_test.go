package test

import (
	"encoding/json"
	"fmt"
	"github.com/wppzxc/taobao_links/pkg/features/coolq/types"
	"golang.org/x/net/websocket"
	"testing"
)

func ConnectWebSocket(wsUrl string) (*websocket.Conn, error) {
	return websocket.Dial(wsUrl, "", "http://api.wpp.pro/")
}

func TestWebSocket(t *testing.T) {
	wsUrl := "ws://49.235.183.87:6700/event/"
	ws, err := ConnectWebSocket(wsUrl)
	if err != nil {
		fmt.Printf("链接 websocket 失败 ： %s，请检查url是否正确！", err)
	}

	for {
		var data []byte
		err := websocket.Message.Receive(ws, &data)
		if err != nil {
			fmt.Println("Error in connect to websocket : ", err)
			continue
		}
		fmt.Printf("get data : %s \n", string(data))
		msg := &types.Message{}
		if err := json.Unmarshal(data, msg); err != nil {
			fmt.Printf("Error in resolve receive data : %s", err)
		}
		fmt.Printf("msg is %#v \n", msg)
	}
}