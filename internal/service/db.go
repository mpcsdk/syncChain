// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"syncChain/internal/model"
	"syncChain/internal/model/entity"
)

type (
	IDB interface {
		Insert(ctx context.Context, data *entity.ChainData) error
		Query(ctx context.Context, query *model.QueryTx) ([]*entity.ChainData, error)
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
