package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"time"
)

var message = make(chan string)

var wsMap = make(map[string]*websocket.Conn)

func upper(ws *websocket.Conn) {
	var err error
	key := time.Now().String()
	wsMap[key] = ws
	for {
		select {
		case msg := <-message:
			suc := 0
			allClient := len(wsMap)
			for k, w := range wsMap {
				if err = websocket.Message.Send(w, msg); err != nil {
					fmt.Println(err)
					delete(wsMap, k)
					w.Close()
					continue
				}
				suc++
			}
			fmt.Printf("send msg %s to %d/%d client, \n", msg, suc, allClient)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	w.Write([]byte("hello world"))
}

func main() {
	http.Handle("/wechat", websocket.Handler(upper))
	http.HandleFunc("/", index)

	go readLogFile()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readLogFile() {
	for {
		//filename := "D:\\project\\go_path\\src\\github.com\\wpp\\taobao_links\\pkg\\features\\wechat\\server\\" + time.Now().Format("2006-01-02") + ".log"
		filename := time.Now().Format("2006-01-02") + ".log"
		cfg := tail.Config{
			Follow: true,
			Logger: tail.DiscardingLogger,
			//Location: &tail.SeekInfo{
			//	Offset: io.SeekEnd,
			//	Whence: io.SeekEnd,
			//},
		}
		t, err := tail.TailFile(filename, cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		for line := range t.Lines {
			//newFilename := "D:\\project\\go_path\\src\\github.com\\wpp\\taobao_links\\pkg\\features\\wechat\\server\\" + time.Now().Format("2006-01-02") + ".log"
			newFilename := time.Now().Format("2006-01-02") + ".log"
			message <- line.Text
			time.Sleep(100 * time.Millisecond)
			if newFilename != filename {
				fmt.Println("change log file to ", newFilename)
				break
			}
		}
	}
}
