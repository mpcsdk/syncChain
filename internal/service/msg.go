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
	IMsg interface {
		SendEvnetBatch(ctx context.Context, datas []*entity.ChainTransfer)
		SendEvent(ctx context.Context, data *entity.ChainTransfer)
	}
)

var (
	localMsg IMsg
)

func Msg() IMsg {
	if localMsg == nil {
		panic("implement not found for interface IMsg, forgot register?")
	}
	return localMsg
}

func RegisterMsg(i IMsg) {
	localMsg = i
}
