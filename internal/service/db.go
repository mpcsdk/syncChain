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
		QueryTransfer(ctx context.Context, chainId int64, query *mpcdao.QueryData) ([]*entity.ChainTransfer, error)
		InsertTransfer(ctx context.Context, chainId int64, data *entity.ChainTransfer) error
		DelChainBlock(ctx context.Context, chainId int64, block int64) error
		InsertTransferBatch(ctx context.Context, chainId int64, datas []*entity.ChainTransfer) error
		ContractAbi() *mpcdao.RiskCtrlRule
		ChainCfg() *mpcdao.ChainCfg
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
