package chaindata

import (
	"context"
	"syncChain/internal/logic/chaindata/block"
	"syncChain/internal/logic/chaindata/common"
	"syncChain/internal/service"

	"syncChain/internal/conf"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
)

type sChainData struct {
	ctx    context.Context
	cancle context.CancelFunc
}

func new() *sChainData {

	ctx, cancle := context.WithCancel(gctx.GetInitCtx())
	s := &sChainData{
		ctx:    ctx,
		cancle: cancle,
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

		g.Log().Notice(s.ctx, "Sycn mode")
		common.InitConf(conf.Config.Chainini)
		block.Init()
		///
		// gproc.AddSigHandlerShutdown(func(sig os.Signal) {
		// 	g.Log().Warning(s.ctx, "Sig:receive signal:", sig.String())

		// 	block.Close()
		// 	//
		// 	// s.cancle()
		// })
		// go gproc.Listen()
	} else {
		g.Log().Notice(s.ctx, "Api mode")
	}
	///
	return s
}
func init() {
	service.RegisterChainData(new())
}
