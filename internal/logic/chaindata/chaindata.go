package chaindata

import (
	"context"
	"os"
	"syncChain/internal/logic/chaindata/block"
	"syncChain/internal/logic/chaindata/common"
	"syncChain/internal/service"

	"syncChain/internal/conf"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
)

type sChainData struct {
	ctx context.Context
}

func new() *sChainData {
	s := &sChainData{
		ctx: gctx.GetInitCtx(),
	}
	//
	p, err := gcmd.Parse(g.MapStrBool{
		"s,sync": false,
	})
	if err != nil {
		panic(err)
	}
	////
	if p.GetOpt("sync") != nil {

		common.InitConf(conf.Config.Chainini)
		block.Init()
		///
		gproc.AddSigHandlerShutdown(func(sig os.Signal) {
			g.Log().Warning(s.ctx, "Sig:receive signal:", sig.String())

			block.Close()
			//
		})
		go gproc.Listen()
	}
	///
	return s
}
func init() {
	service.RegisterChainData(new())
}
