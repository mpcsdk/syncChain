package db

import (
	"syncChain/internal/conf"
	"syncChain/internal/service"

	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mq"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go/jetstream"
)

type sDB struct {
	jet          jetstream.JetStream
	chainData    *mpcdao.ChainData
	riskCtrlRule *mpcdao.RiskCtrlRule
	// s.riskCtrlRule = mpcdao.NewRiskCtrlRule(nil, 0)
}

func (s *sDB) ChainData() *mpcdao.ChainData {
	return s.chainData
}
func (s *sDB) ContractAbi() *mpcdao.RiskCtrlRule {
	return s.riskCtrlRule
}
func new() *sDB {
	nats := mq.New(conf.Config.Nrpc.NatsUrl)
	jet, err := nats.JetStream()
	if err != nil {
		panic(err)
	}
	_, err = nats.GetUpChainTxStream(conf.Config.Server.MsgSize)
	if err != nil {
		panic(err)
	}
	///
	r := g.Redis()
	_, err = r.Conn(gctx.GetInitCtx())
	if err != nil {
		panic(err)
	}
	///
	return &sDB{
		jet:          jet,
		chainData:    mpcdao.NewChainData(r, conf.Config.Cache.SessionDuration),
		riskCtrlRule: mpcdao.NewRiskCtrlRule(r, conf.Config.Cache.SessionDuration),
	}
}
func init() {

	service.RegisterDB(new())
}
