package utils

import (
	"encoding/json"
	"fmt"
	"github.com/wppzxc/taobao_links/pkg/yituike/types"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetItems(token string, url string) (types.ItemResult, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
	}
	//req.Header.Add("token", token)
	req.Header.Set("token", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	if data[0] != '{' {
		fmt.Printf("Error in get premonitor items, the response data is %s\n", data)
		return types.ItemResult{}, fmt.Errorf("Error in get premonitor items, the response data is %s ", data)
	}
	result := types.ItemResult{}
	if err = json.Unmarshal(data, &result); err != nil {
		return types.ItemResult{}, err
	}
	return result, nil
}

func GetToken(auth *types.AuthInfo) (string, error) {

	data := make(url.Values)
	data["account"] = []string{auth.Username}
	data["password"] = []string{auth.Password}
	resp, err := http.PostForm(auth.Url, data)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return "", fmt.Errorf("Error in get token : %s ", err)
	}
	defer resp.Body.Close()
	d, _ := ioutil.ReadAll(resp.Body)
	if d[0] != '{' {
		fmt.Printf("Error in get Token, the result of token is %s\n", d)
		return "", fmt.Errorf("Error in get Token, the result of token is %s ", d)
	}
	result := new(types.TokenResult)
	err = json.Unmarshal(d, result)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
	}
	if result.Code != 0 {
		return "", fmt.Errorf("Error in get token : %s ", err)
	}
	return result.Data[0].Token, nil
}
