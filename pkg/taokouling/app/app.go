package app

import (
	"fmt"
	"net/http"
	"net/url"
)

func GetKouling(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("Error in get response")
		return "-"
	}
	query := resp.Request.URL.RawQuery
	values, err := url.ParseQuery(query)
	if err != nil {
		fmt.Println("Error in get valuse from location")
		return "-"
	}
	v, ok := values["taowords"]
	if !ok {
		fmt.Println("Error in get towards from urlParams")
		return "-"
	}
	if len(v) == 0 {
		fmt.Println("Error in get kouling from towards")
		return "-"
	}
	fmt.Printf("get kouling : %s", v[0])
	return "￥" + v[0][:len(v[0])-1] + "￥"
}
