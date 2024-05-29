package db

import (
	"context"
	"encoding/json"
	"syncChain/internal/conf"
	"syncChain/internal/service"

	"github.com/lib/pq"
	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/mpcsdk/mpcCommon/mq"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go/jetstream"
)

type sDB struct {
	jet           jetstream.JetStream
	chainTransfer *mpcdao.ChainTransfer
	riskCtrlRule  *mpcdao.RiskCtrlRule
	chainCfg      *mpcdao.ChainCfg
	// s.riskCtrlRule = mpcdao.NewRiskCtrlRule(nil, 0)
}

func (s *sDB) QueryTransfer(ctx context.Context, query *mpcdao.QueryData) ([]*entity.ChainTransfer, error) {
	return s.chainTransfer.Query(ctx, query)
}
func isDuplicateKeyErr(err error) bool {
	gerr := err.(*gerror.Error)
	if cerr, ok := gerr.Cause().(*pq.Error); ok {
		if cerr.Code == "23505" {
			return true
		}
	}
	return false
}
func (s *sDB) InsertTransfer(ctx context.Context, data *entity.ChainTransfer) error {
	err := s.chainTransfer.Insert(ctx, data)
	if err != nil {
		if !isDuplicateKeyErr(err) {
			return err
		}
	}
	////sync tx to mq
	d, _ := json.Marshal(data)
	_, err = s.jet.PublishAsync(mq.JetSub_SyncChainTransfer, d)
	if err != nil {
		g.Log().Error(ctx, "InsertTransfer err:", err)
	}
	///
	return nil
}
func (s *sDB) InsertTransferBatch(ctx context.Context, datas []*entity.ChainTransfer) error {
	err := s.chainTransfer.InsertBatch(ctx, datas)
	if err != nil {
		return err
	}
	////sync tx to mq
	for _, data := range datas {
		d, _ := json.Marshal(data)
		s.jet.PublishAsync(mq.JetSub_SyncChainTransfer, d)
	}
	///
	return nil
}

func (s *sDB) ContractAbi() *mpcdao.RiskCtrlRule {
	return s.riskCtrlRule
}
func (s *sDB) ChainCfg() *mpcdao.ChainCfg {
	return s.chainCfg
}
func new() *sDB {
	nats := mq.New(conf.Config.Nrpc.NatsUrl)
	jet := nats.JetStream()
	_, err := nats.CreateOrUpdateStream(mq.JetStream_SyncChain, []string{mq.JetSub_SyncChain}, conf.Config.Server.MsgSize)
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
		jet:           jet,
		chainTransfer: mpcdao.NewChainTransfer(r, conf.Config.Cache.SessionDuration),
		riskCtrlRule:  mpcdao.NewRiskCtrlRule(r, conf.Config.Cache.SessionDuration),
		chainCfg:      mpcdao.NewChainCfg(),
	}
}
func init() {

	service.RegisterDB(new())
}
