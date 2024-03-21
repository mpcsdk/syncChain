package block

import (
	"strings"
	"syncChain/internal/logic/chaindata/common"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	chains       = "chains"
	chainsHeight = "height"
	contracts    = "contracts"
)

var (
	chainClients map[string]*ethModule
	closed       = false
)

func Init() {
	chainClients = make(map[string]*ethModule, 4)

	chainRpcList := common.GlobalConf.GetStrings(chains)
	if 0 == len(chainRpcList) {
		return
	}

	chainContracts := common.GlobalConf.GetStrings(contracts)
	for chain, rpc := range chainRpcList {
		key := strings.ToLower(chain)
		rpcList := strings.Split(rpc, ",")
		module := ethModule{
			ctx:     gctx.GetInitCtx(),
			rpcList: rpcList,
			logger:  g.Log("blocklog"),
		}
		chainClients[key] = &module

		module.start(key, chainContracts[key])
	}

	logLoop()
}

func Close() {
	closed = true
	for _, module := range chainClients {
		module.close()
	}
}
func ClientState() map[int64]int64 {
	d := map[int64]int64{}
	for _, v := range chainClients {
		d[v.chainId] = v.lastBlock
	}
	return d
}
func logLoop() {

	go func() {
		for range time.Tick(time.Second * 10) {
			if closed {
				return
			}

			for _, module := range chainClients {
				g.Log().Notice(gctx.GetInitCtx(), "blockmodule info:", module.info())
			}
		}
	}()
}
