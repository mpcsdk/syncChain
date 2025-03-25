// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type (
	IDB interface {
		InitChainTransferDB(ctx context.Context, chainId int64) error
		QueryTransfer(ctx context.Context, chainId int64, query *mpcdao.QueryData) ([]*entity.SyncchainChainTransfer, error)
		InsertTransfer(ctx context.Context, chainId int64, data *entity.SyncchainChainTransfer) error
		DelChainBlock(ctx context.Context, chainId int64, block int64) error
		UpTransactionMap(ctx context.Context, chainId int64, datas map[int64][]*entity.SyncchainChainTransfer) error
		UpTransaction(ctx context.Context, chainId int64, datas []*entity.SyncchainChainTransfer) error
		InsertTransferBatch(ctx context.Context, chainId int64, datas []*entity.SyncchainChainTransfer) error
		//	func (s *sDB) InsertTransfer_Transaction(ctx context.Context, chainId int64, datas []*entity.SyncchainChainTransfer) error {
		//		// err := s.chainTransfer.InsertBatch(ctx, datas)
		//		chaindb := s.chainTransfer[chainId]
		//		if chaindb == nil {
		//			return errors.New("no chaindb")
		//		}
		//		err := chaindb.Insert_Transaction(ctx, datas)
		//		if err != nil {
		//			return err
		//		}
		//		///
		//		return nil
		//	}
		UpdateState(ctx context.Context, chainId int64, currentBlock int64) error
		GetState(ctx context.Context, chainId int64) (*entity.SyncchainState, error)
		GetContractAbiBriefs(ctx context.Context, chainId int64) ([]*entity.RiskadminContractabi, error)
	}
)

var (
	localDB IDB
)

func DB() IDB {
	if localDB == nil {
		panic("implement not found for interface IDB, forgot register?")
	}
	return localDB
}

func RegisterDB(i IDB) {
	localDB = i
}
