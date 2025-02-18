// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type (
	IEvnetSender interface {
		SendEvnetBatch(ctx context.Context, datas []*entity.SyncchainChainTransfer)
		SendEvnetBatch_Latest(ctx context.Context, datas []*entity.SyncchainChainTransfer)
		SendEvent(ctx context.Context, data *entity.SyncchainChainTransfer)
	}
)

var (
	localEvnetSender IEvnetSender
)

func EvnetSender() IEvnetSender {
	if localEvnetSender == nil {
		panic("implement not found for interface IEvnetSender, forgot register?")
	}
	return localEvnetSender
}

func RegisterEvnetSender(i IEvnetSender) {
	localEvnetSender = i
}
