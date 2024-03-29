package db

import (
	"context"
	"encoding/json"
	"syncChain/internal/conf"
	"syncChain/internal/model"
	"syncChain/internal/service"

	"github.com/mpcsdk/mpcCommon/mpcdao/dao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/mpcsdk/mpcCommon/mq"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/nats-io/nats.go/jetstream"
)

type sDB struct {
	jet jetstream.JetStream
}

func (s *sDB) Insert(ctx context.Context, data *entity.ChainData) error {
	_, err := dao.ChainData.Ctx(ctx).Insert(data)
	if err != nil {
		return err
	}
	////sync tx to mq
	d, _ := json.Marshal(data)
	s.jet.PublishAsync(mq.JetSub_ChainTx, d)
	///
	return nil
}
func (s *sDB) Query(ctx context.Context, query *model.QueryTx) ([]*entity.ChainData, error) {
	if query.PageSize < 1 || query.Page < 0 {
		return nil, nil
	}
	//
	where := dao.ChainData.Ctx(ctx)
	if query.From != "" {
		where = where.Where(dao.ChainData.Columns().From, query.From)
	}
	if query.To != "" {
		where = where.Where(dao.ChainData.Columns().To, query.To)
	}
	if query.Contract != "" {
		where = where.Where(dao.ChainData.Columns().Contract, query.Contract)
	}
	///time
	if query.StartTime != 0 {
		where = where.WhereGTE(dao.ChainData.Columns().Ts, query.StartTime)
	}
	if query.EndTime != 0 {
		where = where.WhereLTE(dao.ChainData.Columns().Ts, query.EndTime)
	}
	///
	if query.PageSize != 0 {
		where = where.Limit(query.Page*query.PageSize, query.PageSize)
	}
	///
	result, err := where.All()
	if err != nil {
		return nil, err
	}
	data := []*entity.ChainData{}
	err = result.Structs(&data)
	///
	return data, err
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

	return &sDB{
		jet: jet,
	}
}
func init() {

	service.RegisterDB(new())
}
