package history

import (
	"github.com/wppzxc/taobao_links/pkg/yituike/types"
	"testing"
)

func TestUpdateHistoryItems(t *testing.T) {
	items := []types.Item{{
		Id: 001,
		StartTime: 1563270600,
	}, {
		Id: 002,
		StartTime: 1563270601,
	}, {
		Id: 003,
		StartTime: 1563270602,
	}}
	UpdateHistoryItems(items)
}
