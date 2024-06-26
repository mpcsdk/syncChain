package chaindata

import (
	"context"
	"syncChain/internal/logic/chaindata/block"
	"syncChain/internal/service"
	"time"

	"syncChain/internal/conf"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mq"
)

type sChainData struct {
	ctx          context.Context
	cancle       context.CancelFunc
	chainclients map[int64]*block.EthModule
	closed       bool
	nats         *mq.NatsServer
	/////
	riskCtrlRule *mpcdao.RiskCtrlRule
	chainCfg     *mpcdao.ChainCfg
	////
}

func (s *sChainData) Close() {
	s.closed = true
	for _, module := range s.chainclients {
		module.Close()
	}
}
func (s *sChainData) ClientState() map[int64]int64 {
	d := map[int64]int64{}
	for _, v := range s.chainclients {
		d[v.ChainId()] = v.LastBlock()
	}
	return d
}

func (s *sChainData) logLoop() {

	go func() {
		for _, module := range s.chainclients {
			g.Log().Notice(gctx.GetInitCtx(), "blockmodule info:", module.Info())
		}

		for range time.Tick(time.Second * 10) {
			if s.closed {
				return
			}

			for _, module := range s.chainclients {
				g.Log().Notice(gctx.GetInitCtx(), "blockmodule info:", module.Info())
			}
		}
	}()
}

////

// //
func New() *sChainData {
	p, err := gcmd.Parse(g.MapStrBool{
		"s,sync": false,
	})
	if err != nil {
		panic(err)
	}
	//
	if p.GetOpt("sync") == nil {
		return &sChainData{}
	}
	////
	////syncchain
	ctx, cancle := context.WithCancel(gctx.GetInitCtx())
	s := &sChainData{
		ctx:          ctx,
		cancle:       cancle,
		chainclients: map[int64]*block.EthModule{},
		closed:       false,
	}
	///db
	s.riskCtrlRule = mpcdao.NewRiskCtrlRule(nil, 0)
	s.chainCfg = mpcdao.NewChainCfg(nil, 0)

	//jet
	nats := mq.New(conf.Config.Nrpc.NatsUrl)
	_, err = nats.CreateOrUpdateStream(mq.JetStream_SyncChain, []string{mq.JetSub_SyncChain}, conf.Config.Server.MsgSize)
	if err != nil {
		panic(err)
	}
	//natsmq
	nats.Sub_ChainCfg(mq.Sub_ChainCfg, func(data *mq.ChainCfgMsg) error {
		g.Log().Notice(gctx.GetInitCtx(), "chaindata:", data)
		switch data.Opt {
		case mq.OptDelete:
			s.deleteOpt(data.Data)
		case mq.OptUpdate:
			s.updateOpt(data.Data)
		case mq.OptAdd:
			s.addOpt(data.Data)
		}
		return nil
	})
	///contractrule
	nats.Sub_ContractRule(mq.Sub_ContractRule, func(data *mq.ContractRuleMsg) error {
		g.Log().Notice(gctx.GetInitCtx(), "contractrule:", data)
		switch data.Opt {
		case mq.OptDelete:
			s.deleteOptContractRule(data.Data)
		case mq.OptUpdate:
			s.updateOptContractRule(data.Data)
		case mq.OptAdd:
			s.addOptContractRule(data.Data)
		}
		return nil
	})

	//natsmq
	nats.Sub_ChainCfg(mq.Sub_ChainCfg, func(data *mq.ChainCfgMsg) error {
		g.Log().Notice(gctx.GetInitCtx(), "chaindata:", data)
		return nil
	})
	///
	g.Log().Notice(s.ctx, "Sycn mode")
	///
	allcfg, err := s.chainCfg.AllCfg(s.ctx)
	if err != nil {
		panic(err)
	}
	for _, v := range allcfg {
		///init db
		err := service.DB().InitChainDB(s.ctx, v.ChainId)
		if err != nil {
			panic(err)
		}
		s.addOpt(v)
	}
	///
	s.logLoop()
	///

	return s
}

func init() {
	// service.RegisterChainData(new())
}
