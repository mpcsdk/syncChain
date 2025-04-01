// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/mpcsdk/mpcCommon/riskAdminService/riskAdminServiceNats"
)

type (
	IDB interface {
		InitChainTransferDB(ctx context.Context, chainId int64) error
		QueryTransfer(ctx context.Context, chainId int64, query *mpcdao.QueryData) ([]*entity.SyncchainChainTransfer, error)
		InsertTransfer(ctx context.Context, chainId int64, data *entity.SyncchainChainTransfer) error
		TruncateTransfer(ctx context.Context, chainId int64, number int64) error
		DelChainBlock(ctx context.Context, chainId int64, block int64) error
		UpTransactionMap(ctx context.Context, chainId int64, datas map[int64][]*entity.SyncchainChainTransfer) error
		UpTransaction(ctx context.Context, chainId int64, datas []*entity.SyncchainChainTransfer) error
		InsertTransferBatch(ctx context.Context, chainId int64, datas []*entity.SyncchainChainTransfer) error
		UpdateState(ctx context.Context, chainId int64, currentBlock int64) error
		GetState(ctx context.Context, chainId int64) (*entity.SyncchainState, error)
		RiskAdminRepo() *riskAdminServiceNats.RiskAdminRepo
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
