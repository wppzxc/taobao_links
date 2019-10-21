package process

import (
	"fmt"
	"github.com/wpp/taobao_links/pkg/yituike/history"
	"github.com/wpp/taobao_links/pkg/yituike/types"
	"github.com/wpp/taobao_links/pkg/yituike/utils"
	"os"
	"time"
)

const (
	proItemsFile = "proItemsFile"
)

type Processer struct {
	Config *types.Config
	Filter func(config *types.Config, items []types.Item) *types.Item
}

func (p *Processer) StartProcess() {
	historyItems := history.GetHistoryItems(proItemsFile)
	
	defer func() {
		history.UpdateHistoryItems(historyItems)
		if err := history.WriteHistoryItems(historyItems, proItemsFile); err != nil {
			fmt.Printf("Error in update history localfile %s \n", err)
		}
	}()
	
	token, err := utils.GetToken(&p.Config.Auth)
	if err != nil {
		return
	}
	result, err := utils.GetItems(token, p.Config.Fanli.Process.Url)
	if err != nil {
		return
	}
	if result.Count > 0 {
	
		diffItems := utils.GetDiffItems(historyItems, result.Data)
		if len(diffItems) == 0 {
			fmt.Println("no new Items found")
			return
		}
	
		// mark all diffItems to historyItems
		for _, i := range diffItems {
			historyItems = append(historyItems, i)
		}
		
		// get filter items
		sendItem := p.Filter(p.Config, diffItems)
		if sendItem == nil {
			fmt.Println("no item found after filter")
			return
		}

		fmt.Printf("get send item %#v \n", sendItem)
	
		tmpfile, err := utils.SaveImage(sendItem.GoodsImageUrl)
		if err != nil {
			fmt.Println("Error in create tmp image file")
			return
		}
		defer os.Remove(tmpfile.Name())
	
		// send message to users
		for _, u := range p.Config.Receivers {
			if err := utils.SendImage(tmpfile, u); err != nil {
				fmt.Printf("Send image to user %s error : %s\n", u.Name, err)
			}
	
			msg := utils.GetMsg(p.Config.Fanli.Process.MsgPrefix, sendItem, u.Link)
			if err := utils.SendMessage(msg, u.Name); err != nil {
				fmt.Printf("Send message to user %s error : %s\n", u.Name, err)
			}
			time.Sleep(time.Duration(p.Config.Fanli.SendInterval) * time.Second)
		}
	} else {
		fmt.Println("no process items found")
		return
	}
}
