package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
	"k8s.io/klog"
)

type WsConnection struct {
	Message chan *WeChatMessage
	Stop    chan string
}

type ReceiveMessage struct {
	ToUserName string `xml:"ToUserName"`
	AgentID    string `xml:"AgentID"`
	Encrypt    string `xml:"Encrypt"`
}

type WeChatMessage struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgID        string `xml:"MsgId"`
	AgentID      string `xml:"AgentID"`
}

var wsMap = new(sync.Map)

var encodingAesKey = "R6niko99y1tVWd1BlMAhVKY5edZd9Wyo1QUc1vQw6qV"
// var encodingAesKey = "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/clients", allClients)
	e.GET("/receive", receive)
	e.POST("/receive", receive)
	e.GET("/listen/:wxid", listen)
	e.Logger.Fatal(e.Start(":8080"))
}

func allClients(c echo.Context) error {
	wxids := make([]string, 0)
	wsMap.Range(func(key, value interface{}) bool {
		wxids = append(wxids, fmt.Sprintf("%s", key))
		return true
	})
	return c.JSON(http.StatusOK, wxids)
}

func receive(c echo.Context) error {
	msg_signature := c.QueryParam("msg_signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	echostr := c.QueryParam("echostr")
	if len(msg_signature) == 0 || len(timestamp) == 0 || len(nonce) == 0 {
		return c.JSON(http.StatusBadRequest, "QueryParam can't be null")
	}

	if len(echostr) > 0 {
		originStr, err := base64.StdEncoding.DecodeString(echostr)
		if err != nil {
			klog.Error(err)
			return err
		}
		msg, err := testDecodeMsg(originStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.String(http.StatusOK, msg)

	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	klog.Info(string(data))

	toUser, wm, receiveid, err := decodeMsg(data)
	if err != nil {
		klog.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	klog.Infof("msg_signature : '%s', timestamp : '%s', nonce : '%s', msg : '%s', receiveid : '%s'", msg_signature, timestamp, nonce, wm.Content, receiveid)
	wsMap.Range(func(key, value interface{}) bool {
		if key == toUser {
			conn := value.(WsConnection)
			conn.Message <- wm
			return false
		}
		return true
	})

	return nil
}

func listen(c echo.Context) error {
	wxid := c.Param("wxid")
	if conn, ok := wsMap.Load(wxid); ok {
		if c, ok := conn.(WsConnection); ok {
			c.Stop <- "stop by new conn"
		}
	}
	wsConn := WsConnection{
		Message: make(chan *WeChatMessage),
		Stop:    make(chan string),
	}
	wsMap.Store(wxid, wsConn)

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		defer wsMap.Delete(wxid)

		for {
			select {
			case msg := <-wsConn.Message:
				data, err := json.Marshal(msg)
				if err != nil {
					klog.Error(err)
					return
				}
				if err := websocket.Message.Send(ws, string(data)); err != nil {
					klog.Error(err)
					return
				}
				klog.Infof("send msg %s to wechat: %s", msg, wxid)
			case stop := <-wsConn.Stop:
				klog.Warningf("websocket server for wxid '%s' stop, because '%s'", wxid, stop)
				wsMap.Delete(wxid)
				ws.Close()
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

func testDecodeMsg(data []byte) (string, error) {

	aesKey, err := base64.StdEncoding.DecodeString(encodingAesKey + "=")
	if err != nil {
		klog.Error(err)
		return "", err
	}
	randMsg := AesDecryptCBC(data, aesKey)

	// runeMsg := []rune(string(randMsg))
	content := randMsg[16:]
	msgLen := content[:4]
	length, err := strconv.ParseUint(hex.EncodeToString(msgLen), 16, 32)
	if err != nil {
		klog.Error(err)
		return "", err
	}

	klog.Info(msgLen)
	klog.Info(length)
	msg := content[4 : length+4]
	return string(msg), nil
}

func decodeMsg(data []byte) (string, *WeChatMessage, string, error) {
	rm := new(ReceiveMessage)
	if err := xml.Unmarshal(data, rm); err != nil {
		klog.Error(err)
		return "", nil, "", err
	}

	aesKey, err := base64.StdEncoding.DecodeString(encodingAesKey + "=")
	if err != nil {
		klog.Error(err)
		return "", nil, "", err
	}

	aesMsg, err := base64.StdEncoding.DecodeString(rm.Encrypt)
	if err != nil {
		klog.Error(err)
		return "", nil, "", err
	}
	randMsg := AesDecryptCBC(aesMsg, aesKey)

	content := randMsg[16:]
	msgLen := content[:4]
	length, err := strconv.ParseUint(hex.EncodeToString(msgLen), 16, 32)
	if err != nil {
		klog.Error(err)
		return "", nil, "", err
	}
	msg := content[4 : length+4]
	receiveid := content[len(msg)+4:]

	wm := new(WeChatMessage)
	if err := xml.Unmarshal([]byte(msg), wm); err != nil {
		klog.Error(err)
		return "", nil, "", err
	}

	return rm.ToUserName, wm, string(receiveid), nil
}

func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}
func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
