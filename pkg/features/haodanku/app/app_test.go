package app

import (
	"fmt"
	"testing"
)

func TestGetItems(t *testing.T) {
	url := "https://www.haodanku.com/indexapi/get_allitem_list?sort=3&p=1"
	items, err := GetLinks(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(items)
}
