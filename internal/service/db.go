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
		QueryTransfer(ctx context.Context, query *mpcdao.QueryData) ([]*entity.ChainTransfer, error)
		InsertTransfer(ctx context.Context, data *entity.ChainTransfer) error
		InsertTransferBatch(ctx context.Context, datas []*entity.ChainTransfer) error
		ContractAbi() *mpcdao.RiskCtrlRule
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
