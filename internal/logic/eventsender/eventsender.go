package msg

import (
	"context"
	"encoding/json"
	"syncChain/internal/conf"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/mpcsdk/mpcCommon/mq"
	"github.com/nats-io/nats.go/jetstream"
)

type sEvnetSender struct {
	jet jetstream.JetStream
}

func NewMsg() *sEvnetSender {
	nats := mq.New(conf.Config.Nrpc.NatsUrl)
	jet := nats.JetStream()
	_, err := nats.CreateOrUpdateStream(mq.JetStream_SyncChain, []string{mq.JetSub_SyncChain}, conf.Config.Syncing.MsgSize)
	if err != nil {
		panic(err)
	}

	///
	return &sEvnetSender{
		jet: jet,
	}

}
func (s *sEvnetSender) SendEvnetBatch(ctx context.Context, datas []*entity.SyncchainChainTransfer) {
	////sync tx to mq
	for _, data := range datas {
		d, _ := json.Marshal(data)
		_, err := s.jet.PublishAsync(mq.JetSub_SyncChainTransfer, d)
		if err != nil {
			g.Log().Error(ctx, "SendMsg err:", err, "data:", data)
		}
	}
	///
}
func (s *sEvnetSender) SendEvnetBatch_Latest(ctx context.Context, datas []*entity.SyncchainChainTransfer) {
	////sync tx to mq
	for _, data := range datas {
		d, _ := json.Marshal(data)
		_, err := s.jet.PublishAsync(mq.JetSub_SyncChainTransfer_Latest, d)
		if err != nil {
			g.Log().Error(ctx, "SendMsg err:", err, "data:", data)
		}
	}
	///
}
func (s *sEvnetSender) SendEvent(ctx context.Context, data *entity.SyncchainChainTransfer) {
	////sync tx to mq
	d, _ := json.Marshal(data)
	_, err := s.jet.PublishAsync(mq.JetSub_SyncChainTransfer, d)
	if err != nil {
		g.Log().Error(ctx, "SendMsg err:", err)
	}
	///
}
