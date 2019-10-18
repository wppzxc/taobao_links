package app

import (
	"github.com/wpp/taobao_links/pkg/yituike/types"
	"github.com/wpp/taobao_links/pkg/yituike/wait"
	"time"
)

func Start(conf *types.Config, fs []func(), stopCh chan struct{}) {
	for _, f := range fs {
		go wait.Until(f, time.Duration(conf.Fanli.RefreshInterval)*time.Second, stopCh)
	}
}

