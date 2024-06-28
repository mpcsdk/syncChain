package transfer

import (
	"context"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func ProcessInTxnsRpg(ctx context.Context, chainId int64, block *types.Block, traces []*util.TraceRpg) []*entity.ChainTransfer {

	////
	filtertrace := []*util.TraceRpg{}
	for _, trace := range traces {
		if trace.Type == "call" {
			if trace.Value != "0" {
				filtertrace = append(filtertrace, trace)
			}
		}
	}

	//// fill transfer
	transfers := []*entity.ChainTransfer{}
	for _, trace := range filtertrace {
		tx := block.Transaction(trace.ParentTxHash)
		if tx == nil {
			g.Log().Warning(ctx, "tx is nil")
			continue
		}
		/////
		transfer := &entity.ChainTransfer{
			ChainId:   chainId,
			Height:    trace.BlockHeight,
			BlockHash: trace.BlockHash.String(),
			Ts:        int64(block.Time()),
			TxHash:    trace.ParentTxHash.String(),
			TxIdx:     trace.TxIndex,
			From:      trace.Source.String(),
			To:        trace.Target.String(),
			Contract:  "",
			Value:     trace.Value,
			Gas:       trace.GasLimit,
			GasPrice:  "0",
			LogIdx:    -1,
			Nonce:     int64(tx.Nonce()),
			Kind:      "external",
			Status:    0,
			Removed:   false,
			TraceTag:  trace.TraceTag,
		}
		transfers = append(transfers, transfer)
	}

	return transfers
}
func ProcessInTxns(ctx context.Context, chainId int64, block *types.Block, traces []*util.Trace) []*entity.ChainTransfer {
	///

	////
	filtertrace := []*util.Trace{}
	for _, trace := range traces {
		if trace.Action.CallType == "call" {
			if trace.Action.Value != "0x0" {
				filtertrace = append(filtertrace, trace)
			}
		}
	}

	//// fill transfer
	transfers := []*entity.ChainTransfer{}
	for _, trace := range filtertrace {
		tx := block.Transaction(trace.TransactionHash)
		if tx == nil {
			g.Log().Warning(ctx, "tx is nil")
			continue
		}
		/////
		transfer := &entity.ChainTransfer{
			ChainId:   chainId,
			Height:    trace.BlockNumber,
			BlockHash: trace.BlockHash.String(),
			Ts:        int64(block.Time()),
			TxHash:    trace.TransactionHash.String(),
			TxIdx:     trace.TransactionPosition,
			From:      trace.Action.From.String(),
			To:        trace.Action.To.String(),
			Contract:  "",
			Value:     trace.Action.Value,
			Gas:       trace.Action.Gas,
			GasPrice:  "0",
			LogIdx:    -1,
			Nonce:     int64(tx.Nonce()),
			Kind:      "external",
			Status:    0,
			Removed:   false,
			TraceTag:  trace.Tag(),
		}
		transfers = append(transfers, transfer)
	}

	return transfers
}
