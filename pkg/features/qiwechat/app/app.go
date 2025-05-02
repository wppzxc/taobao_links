package app

import (
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	SendTime string
	Receiver string
	MsgType  string
	Sender   string
	Text     string
}

func ConnectWebSocket(wsUrl string) (*websocket.Conn, error) {
	return websocket.Dial(wsUrl, "", "http://api.wpp.pro/")
}

func Start(wsUrl string, groups []string, users []string, interval int, tklTitle string, filterNo []string, filterYes []string, stopCh chan struct{}) error {
	// start websocket client goroutine
	msgs := make(chan Message, 10)
	ws, err := ConnectWebSocket(wsUrl)
	if err != nil {
		return fmt.Errorf("链接 websocket 失败 ： %s，请检查url是否正确！", err)
	}
	go func(ws *websocket.Conn, msgs chan Message, stopCh chan struct{}) {
		defer ws.Close()
		tmpMessage := ""
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
				msg := string(data)
				// if msg = ""
				if len(msg) == 0 {
					tmpMessage = fmt.Sprintf("%s\n%s", tmpMessage, msg)
					continue
				}
				// if msg = "**** \r", send the message
				if msg[len(msg)-1:] == "\r" {
					tmpMessage = fmt.Sprintf("%s\n%s", tmpMessage, msg[:len(msg)-1])
					message := newMessage(tmpMessage)
					if message == nil {
						tmpMessage = ""
						continue
					}
					// 发送时间超过1分钟，则丢弃
					sendTime, _ := time.ParseInLocation("2006-01-02 15:04:05", strings.Trim(message.SendTime, "\n"), time.Local)
					if sendTime.Add(1 * time.Minute).After(time.Now()) {
						continue
					}

					// 过滤系统消息
					if message == nil {
						continue
					}
					// 屏蔽关键词
					if skip := filterWechatMessage(message.Text, filterNo); skip {
						tmpMessage = ""
						continue
					}

					// 选择关键词
					if len(filterYes) > 0 {
						if isSend := filterWechatMessage(message.Text, filterYes); isSend {
							if ok := checkSendGroups(message.Sender, groups); ok {
								msgs <- *message
							}
						}
						tmpMessage = ""
						// 没有关键词则全部转发
					} else {
						if ok := checkSendGroups(message.Sender, groups); ok {
							msgs <- *message
						}
						tmpMessage = ""
					}
				} else {
					tmpMessage = fmt.Sprintf("%s\n%s", tmpMessage, msg)
					continue
				}
			}
		}
	}(ws, msgs, stopCh)

	// start sender goroutine
	go func(msgs chan Message, users []string, interval int, tklTitle string, filterNo []string, filterYes []string, stopCh chan struct{}) {
		for {
			select {
			case msg := <-msgs:
				if err := WechatMessageSend(msg, users, interval, tklTitle); err != nil {
					fmt.Println("Error in send meg : ", err)
				} else {
					fmt.Printf("sender message : %#v \n", msg)
				}
			case <-stopCh:
				fmt.Println("sender stoped!")
				return
			}
		}
	}(msgs, users, interval, tklTitle, filterNo, filterYes, stopCh)
	return nil
}

func checkSendGroups(sender string, groups []string) bool {
	// groupId := strconv.FormatInt(msg.GroupId, 10)
	for _, g := range groups {
		if strings.Index(sender, g) >= 0 {
			return true
		}
	}
	return false
}

func CheckCoolqLogined() bool {
	_, err := net.DialTimeout("tcp", "0.0.0.0:6700", 3*time.Second)
	if err != nil {
		return false
	}
	return true
}

func newMessage(msg string) *Message {
	strs := strings.Split(msg, "->")
	if len(strs) < 4 {
		return nil
	}
	sender, text := getMessageSenderAndText(strs[3])
	message := &Message{
		SendTime: strs[0][0 : len(strs[0])-1],
		Receiver: strs[1],
		MsgType:  strs[2],
		Sender:   sender,
		Text:     text,
	}
	return message
}

func getMessageSenderAndText(str string) (string, string) {
	strs := strings.SplitN(str, "说:", 2)
	if len(strs) == 2 {
		return strs[0], strs[1]
	}
	return "", ""
}

func filterWechatMessage(msg string, filters []string) bool {
	for _, f := range filters {
		if strings.Index(msg, f) >= 0 {
			return true
		}
	}
	return false
}
