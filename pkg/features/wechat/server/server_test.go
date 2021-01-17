package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	wsUrl := "ws://81.70.210.167:8080/listen/ww86f1fcc69789b25b"
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
	wsUrl := "ws://127.0.0.1:8080/listen/wx5823bf96d3bd56c7a"
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
