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
)

type sDB struct {
	r             *gredis.Redis
	dur           int
	chainTransfer map[int64]*mpcdao.ChainTransfer
	riskCtrlRule  *mpcdao.RiskCtrlRule
	chainCfg      *mpcdao.ChainCfg
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
func (s *sDB) InitChainDB(ctx context.Context, chainId int64) error {
	err := mpcdao.CreateChainDB(ctx, chainId)
	if err != nil {
		if isPgErr(err, "42P04") {
		} else {
			return err
		}
	}
	chaindb := mpcdao.NewChainTransfer(chainId, s.r, s.dur)
	s.chainTransfer[chainId] = chaindb
	return nil
}
func (s *sDB) QueryTransfer(ctx context.Context, chainId int64, query *mpcdao.QueryData) ([]*entity.ChainTransfer, error) {
	// return s.chainTransfer.Query(ctx, query)
	if chaindb, ok := s.chainTransfer[chainId]; ok {
		return chaindb.Query(ctx, query)
	} else {
		return nil, errors.New("no chaindb")
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
func (s *sDB) InsertTransfer(ctx context.Context, chainId int64, data *entity.ChainTransfer) error {
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
func (s *sDB) DelChainBlock(ctx context.Context, chainId int64, block int64) error {
	chaindb := s.chainTransfer[chainId]
	if chaindb == nil {
		return errors.New("no chaindb")
	}
	err := chaindb.DelChainBlockNumber(ctx, chainId, block)
	return err

}
func (s *sDB) InsertTransferBatch(ctx context.Context, chainId int64, datas []*entity.ChainTransfer) error {
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

func (s *sDB) ContractAbi() *mpcdao.RiskCtrlRule {
	return s.riskCtrlRule
}
func (s *sDB) ChainCfg() *mpcdao.ChainCfg {
	return s.chainCfg
}
func New() *sDB {

	///
	r := g.Redis()
	_, err := r.Conn(gctx.GetInitCtx())
	if err != nil {
		panic(err)
	}
	///
	return &sDB{
		r:             r,
		dur:           conf.Config.Cache.SessionDuration,
		chainTransfer: map[int64]*mpcdao.ChainTransfer{},
		//mapmpcdao.NewChainTransfer(r, conf.Config.Cache.SessionDuration),
		riskCtrlRule: mpcdao.NewRiskCtrlRule(r, conf.Config.Cache.SessionDuration),
		chainCfg:     mpcdao.NewChainCfg(r, conf.Config.Cache.SessionDuration),
	}
}
func init() {

}
