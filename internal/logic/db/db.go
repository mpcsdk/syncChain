package db

import (
	"context"
	"encoding/json"
	"syncChain/internal/conf"
	"syncChain/internal/service"

	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/mpcsdk/mpcCommon/mq"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go/jetstream"
)

type sDB struct {
	jet           jetstream.JetStream
	chainTransfer *mpcdao.ChainTransfer
	riskCtrlRule  *mpcdao.RiskCtrlRule
	// s.riskCtrlRule = mpcdao.NewRiskCtrlRule(nil, 0)
}

func (s *sDB) QueryTransfer(ctx context.Context, query *mpcdao.QueryData) ([]*entity.ChainTransfer, error) {
	return s.chainTransfer.Query(ctx, query)
}

func (s *sDB) InsertTransfer(ctx context.Context, data *entity.ChainTransfer) error {
	err := s.chainTransfer.Insert(ctx, data)
	if err != nil {
		return err
	}
	////sync tx to mq
	d, _ := json.Marshal(data)
	s.jet.PublishAsync(mq.JetSub_SyncChainTransfer, d)
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
	}
}
func init() {

	service.RegisterDB(new())
}
