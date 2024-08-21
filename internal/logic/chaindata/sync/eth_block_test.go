package syncBlock

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestBlock(t *testing.T) {
	ctx := gctx.GetInitCtx()
	module := NewEthModule(
		ctx,
		"RPG",
		2025,
		67709977,
		[]string{
			"https://mainnet-rpc.rangersprotocol.com/api/jsonrpc",
		},
		[]common.Address{},
		[]common.Address{},
	)
	module.Start()
	time.Sleep(10 * time.Minute)
}
