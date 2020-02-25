package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	wsUrl := "ws://49.233.131.124:8080/wechat"
	ws, err := websocket.Dial(wsUrl, "", "http://api.wpp.pro1/")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for {
		var data []byte
		err := websocket.Message.Receive(ws, &data)
		if err != nil {
			fmt.Println("Error in get message : ", err)
			continue
		}
		fmt.Printf("get message is : %#v \n", string(data))
	}
}

func TestClient2(t *testing.T) {
	wsUrl := "ws://127.0.0.1:8080/wechat"
	ws, err := websocket.Dial(wsUrl, "", "http://api.wpp.pro2/")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for {
		var data []byte
		err := websocket.Message.Receive(ws, &data)
		if err != nil {
			fmt.Println("Error in get message : ", err)
			continue
		}
		fmt.Printf("get message is : %#v \n", string(data))
	}
}
