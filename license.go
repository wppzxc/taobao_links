package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	key = "key-v1.x"
	licenseFile = "license.txt"
)

var Features = []string{
	"dataoke",
	"haodanku",
	"duoduojinbao",
	"taokeyi",
	"goodSearch",
	"yituike",
	"pddUserNumber",
	"taokouling",
	"coolq",
}

type License struct {
	Feature         []string `json:"feature"`
	ExpireTimestamp int64    `json:"expireTimestamp"`
}

func GetLocalLicense() *License {
	_, err := os.Stat(licenseFile)
	if os.IsNotExist(err) {
		_, err := os.Create(licenseFile)
		if err != nil {
			fmt.Printf("Error : %s\n", err)
		}
		return nil
	}
	data, err := ioutil.ReadFile(licenseFile)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	decodeData := DecodeLicense(string(data))
	if len(decodeData) == 0 {
		return nil
	}
	license := &License{}
	if err := json.Unmarshal([]byte(decodeData), license); err != nil {
		fmt.Printf("Error : %s\n", err)
		return nil
	}
	return license
}

func CheckLicense(licenseStr string) error {
	license := &License{}
	decodeData := DecodeLicense(licenseStr)
	if err := json.Unmarshal([]byte(decodeData), license); err != nil {
		return fmt.Errorf("Error in decode license : %s\n", err)
	}
	now := time.Now().Unix()
	if license.ExpireTimestamp > now {
		if err := ioutil.WriteFile(licenseFile, []byte(licenseStr), 0755); err != nil {
			return fmt.Errorf("Error in save license : %s\n", err)
		}
		return nil
	}
	return fmt.Errorf("your license are expired! ")
}

func ResetLocalLicense() error {
	if err := ioutil.WriteFile(licenseFile, []byte(""), 0755); err != nil {
		return err
	}
	return nil
}

/* license
   encode by data -> des -> base64 -> encodeData
    decode by encodeData -> base64 -> des
*/
// 加密license
func DecodeLicense(encodeData string) string {
	// decode base64
	decodeBase64Data, err := base64.StdEncoding.DecodeString(encodeData)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// decode des
	resultData := MyDESDecrypt(string(decodeBase64Data), []byte(key))
	return resultData
}

/* license
encode by data -> des -> base64 -> base64 -> encodeData
 decode by encodeData -> base64 -> base64 -> des
*/
// 解密license
func EncodeLicense(license *License) string {
	// get data
	data, _ := json.Marshal(license)
	// get des and base64 data
	desData := MyDESEncrypt(data, []byte(key))
	// get base64 data
	base64Data := base64.StdEncoding.EncodeToString([]byte(desData))
	return base64Data
}

//DES加密方法
func MyDESEncrypt(origData, key []byte) string {
	//将字节秘钥转换成block快
	block, _ := des.NewCipher(key)
	//对明文先进行补码操作
	origData = PKCS5Padding(origData, block.BlockSize())
	//设置加密方式
	blockMode := cipher.NewCBCEncrypter(block, key)
	//创建明文长度的字节数组
	crypted := make([]byte, len(origData))
	//加密明文,加密后的数据放到数组中
	blockMode.CryptBlocks(crypted, origData)
	//将字节数组转换成字符串
	return base64.StdEncoding.EncodeToString(crypted)
}

//实现明文的补码
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	//计算出需要补多少位
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把参数一 切片复制 参数二count个,然后合成一个新的字节切片返回
	// 需要补padding位的padding值
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	//把补充的内容拼接到明文后面
	return append(ciphertext, padtext...)
}

//解密
func MyDESDecrypt(data string, key []byte) string {
	//倒叙执行一遍加密方法
	//将字符串转换成字节数组
	crypted, _ := base64.StdEncoding.DecodeString(data)
	//将字节秘钥转换成block快
	block, _ := des.NewCipher(key)
	//设置解密方式
	blockMode := cipher.NewCBCDecrypter(block, key)
	//创建密文大小的数组变量
	origData := make([]byte, len(crypted))
	//解密密文到数组origData中
	blockMode.CryptBlocks(origData, crypted)
	//去补码
	origData = PKCS5UnPadding(origData)

	return string(origData)
}

//去除补码
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	//解密去补码时需取最后一个字节，值为m，则从数据尾部删除m个字节，剩余数据即为加密前的原文
	return origData[:(length - unpadding)]
}
