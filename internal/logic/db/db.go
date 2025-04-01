package db

import (
	"context"
	"errors"
	"syncChain/internal/conf"

	"github.com/lib/pq"
	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/mpcsdk/mpcCommon/riskAdminService/riskAdminServiceNats"
)

type sDB struct {
	r                *gredis.Redis
	dur              int
	chainTransfer    map[int64]*mpcdao.ChainTransfer
	riskAdminService *riskAdminServiceNats.RiskAdminNatsService
	// chainCfg      *mpcdao.ChainCfg
}

func isPgErr(err error, key string) bool {
	gerr := err.(*gerror.Error)
	if cerr, ok := gerr.Cause().(*pq.Error); ok {
		if cerr.Code == pq.ErrorCode(key) {
			return true
		}
	}
	return false
}
func (s *sDB) InitChainTransferDB(ctx context.Context, chainId int64) error {
	err := mpcdao.InitSyncChainDB(ctx, chainId)
	if err != nil {
		panic(err)
	}
	chaindb := mpcdao.NewChainTransfer(chainId, s.r, s.dur)
	s.chainTransfer[chainId] = chaindb
	return nil
}
func (s *sDB) QueryTransfer(ctx context.Context, chainId int64, query *mpcdao.QueryData) ([]*entity.SyncchainChainTransfer, error) {
	// return s.chainTransfer.Query(ctx, query)
	if chaindb, ok := s.chainTransfer[chainId]; ok {
		return chaindb.Query(ctx, query)
	} else {
		g.Log().Error(ctx, "QueryTransfer:", "chainId:", chainId, "query:", query)
		return nil, nil
	}
}

// /
func isDuplicateKeyErr(err error) bool {
	gerr := err.(*gerror.Error)
	if cerr, ok := gerr.Cause().(*pq.Error); ok {
		if cerr.Code == "23505" {
			return true
		}
	}
	return false
}
func (s *sDB) InsertTransfer(ctx context.Context, chainId int64, data *entity.SyncchainChainTransfer) error {
	// err := s.chainTransfer.Insert(ctx, data)
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return errors.New("no chaindb")
	}

	err := chaindb.Insert(ctx, data)
	if err != nil {
		if !isDuplicateKeyErr(err) {
			return err
		}
	}

	return nil
}
func (s *sDB) TruncateTransfer(ctx context.Context, chainId int64, number int64) error {
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return errors.New("no chaindb")
	}
	err := chaindb.TruncateTransfer(ctx, chainId, number)
	return err
}
func (s *sDB) DelChainBlock(ctx context.Context, chainId int64, block int64) error {
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return errors.New("no chaindb")
	}
	err := chaindb.DelChainBlockNumber(ctx, chainId, block)
	return err
}
func (s *sDB) UpTransactionMap(ctx context.Context, chainId int64, datas map[int64][]*entity.SyncchainChainTransfer) error {
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return errors.New("no chaindb")
	}
	err := chaindb.UpTransactionMap(ctx, datas)
	if err != nil {
		return err
	}
	return nil
}

func (s *sDB) UpTransaction(ctx context.Context, chainId int64, datas []*entity.SyncchainChainTransfer) error {
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return errors.New("no chaindb")
	}
	err := chaindb.UpTransaction(ctx, datas)
	if err != nil {
		return err
	}
	return nil
}
func (s *sDB) InsertTransferBatch(ctx context.Context, chainId int64, datas []*entity.SyncchainChainTransfer) error {
	// err := s.chainTransfer.InsertBatch(ctx, datas)
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return errors.New("no chaindb")
	}
	err := chaindb.InsertBatch(ctx, datas)
	if err != nil {
		return err
	}
	///
	return nil
}

func (s *sDB) UpdateState(ctx context.Context, chainId int64, currentBlock int64) error {
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return errors.New("no chaindb")
	}
	err := chaindb.UpdateState(ctx, chainId, currentBlock)
	if err != nil {
		return err
	}
	///
	return nil

}
func (s *sDB) GetState(ctx context.Context, chainId int64) (*entity.SyncchainState, error) {
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return nil, errors.New("no chaindb")
	}
	return chaindb.GetState(ctx, chainId)
}

func (s *sDB) RiskAdminRepo() *riskAdminServiceNats.RiskAdminRepo {
	return s.riskAdminService.RiskAdminRepo()
}
func New() *sDB {

	///
	r := g.Redis()
	_, err := r.Conn(gctx.GetInitCtx())
	if err != nil {
		panic(err)
	}
	////
	///
	riskAdminService, err := riskAdminServiceNats.NewRiskAdminNatsService(gctx.GetInitCtx(),
		riskAdminServiceNats.RiskAdminServiceCfgCfgBuilder().
			// WithConsumeChainFn(s.ConsumeChainCfg).
			// WithConsumeContractFn(s.ConsumeContract).
			// WithConsumeRiskRuleFn(s.ConsumeRiskRule).
			// WithConsumeRiskRuleCheckRespFn(s.consumeRiskRuleCheck).
			WithRedis(r, conf.Config.Cache.Duration).
			WithUrlTimeOut(conf.Config.Nats.NatsUrl, int64(conf.Config.Nats.TimeOut)),
	)
	if err != nil {
		panic(err)
	}
	///
	s := &sDB{
		r: r,
		// dur:           conf.Config.Cache.SessionDuration,
		chainTransfer:    map[int64]*mpcdao.ChainTransfer{},
		riskAdminService: riskAdminService,
		//mapmpcdao.NewChainTransfer(r, conf.Config.Cache.SessionDuration),
		// riskCtrlRule: mpcdao.NewRiskCtrlRule(r, conf.Config.Cache.SessionDuration),
		// chainCfg:     mpcdao.NewChainCfg(r, conf.Config.Cache.SessionDuration),
	}
	return s
}
func init() {

}
