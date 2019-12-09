package app

import (
	"github.com/wppzxc/taobao_links/pkg/features/yituike/types"
	"github.com/wppzxc/taobao_links/pkg/features/yituike/wait"
	"time"
)

func Start(conf *types.Config, fs []func(), stopCh chan struct{}) {
	for _, f := range fs {
		go wait.Until(f, time.Duration(conf.Fanli.RefreshInterval)*time.Second, stopCh)
	}
}

