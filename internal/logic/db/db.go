package db

import (
	"context"
	"encoding/json"
	"syncChain/internal/conf"
	"syncChain/internal/service"

	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/mpcsdk/mpcCommon/mq"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go/jetstream"
)

type sDB struct {
	jet       jetstream.JetStream
	chainData *mpcdao.ChainData
}

func (s *sDB) Insert(ctx context.Context, data *entity.ChainData) error {
	err := s.chainData.Insert(ctx, data)
	if err != nil {
		return err
	}
	////sync tx to mq
	d, _ := json.Marshal(data)
	s.jet.PublishAsync(mq.JetSub_ChainTx, d)
	///
	return nil
}
func (s *sDB) Query(ctx context.Context, query *mpcdao.QueryData) ([]*entity.ChainData, error) {
	return s.chainData.Query(ctx, query)
}

func (s *sDB) ChainData() *mpcdao.ChainData {
	return s.chainData
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
		jet:       jet,
		chainData: mpcdao.NewChainData(r, conf.Config.Cache.SessionDuration),
	}
}
func init() {

	service.RegisterDB(new())
}
