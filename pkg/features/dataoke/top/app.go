package top

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/wppzxc/taobao_links/pkg/features/dataoke/types"
	"github.com/wppzxc/taobao_links/pkg/utils/dataokeapi"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	openapiGetGoodsListHost   = "https://openapi.dataoke.com/api/goods/get-goods-list"
	openapiGetRankingListHost = "https://openapi.dataoke.com/api/goods/get-ranking-list"
	openapiGetGoodsDetailHost = "https://openapi.dataoke.com/api/goods/get-goods-details"
	appKey                    = "5e9d2dbadc286"
	appSecret                 = "8f3c81484fdf7bd2695ddbbc6a128201"
	apiVersion                = "v1.2.2"
)

func GetTblinksByApi(goodsIds []string) string {
	text := ""
	for _, id := range goodsIds {
		goods, err := GetGoodsDetailApi(id)
		if err != nil {
			continue
		}
		text = text + "\n" + goods.ItemLink
	}
	return text
}

func GetGoodsDetailApi(goodsId string) (*types.DataokeData, error) {
	param := url.Values{
		"appKey":   []string{appKey},
		"version":  []string{apiVersion},
		"goodsId": []string{goodsId},
	}
	sign := dataokeapi.MakeSign(param.Encode(), appSecret)
	param["sign"] = []string{sign}
	reqUrl := openapiGetGoodsDetailHost + "?" + param.Encode()
	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respData := new(types.DataokeGoodsDetailResponseData)
	if err := json.Unmarshal(data, respData); err != nil {
		return nil, err
	}
	return &respData.Data, nil
}

func GetRankingListApi() ([]types.DataokeData, error) {
	param := url.Values{
		"appKey":   []string{appKey},
		"version":  []string{apiVersion},
		"rankType": []string{"1"},
	}
	sign := dataokeapi.MakeSign(param.Encode(), appSecret)
	param["sign"] = []string{sign}
	reqUrl := openapiGetRankingListHost + "?" + param.Encode()
	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respData := new(types.DataokeTopApiResponseData)
	if err := json.Unmarshal(data, respData); err != nil {
		return nil, err
	}
	return respData.Data, nil
}

func GetGoodsListApi(page string, cid string) (*types.DataokeQuanDataList, error) {
	param := url.Values{
		"appKey":   []string{appKey},
		"version":  []string{apiVersion},
		"pageSize": []string{"100"},
		"pageId":   []string{page},
		"sort":     []string{"2"},
	}
	if len(cid) > 0 {
		param["cids"] = []string{cid}
	}
	sign := dataokeapi.MakeSign(param.Encode(), appSecret)
	param["sign"] = []string{sign}
	reqUrl := openapiGetGoodsListHost + "?" + param.Encode()
	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respData := new(types.DataokeQuanApiResponseData)
	if err := json.Unmarshal(data, respData); err != nil {
		return nil, err
	}
	return &respData.Data, nil
}

func makeSign(param string) string {
	str := param + "&key=" + appSecret
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
