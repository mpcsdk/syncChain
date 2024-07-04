package transfer

import (
	"context"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func ProcessInTxnsRpg(ctx context.Context, chainId int64, block *types.Block, traces []*util.TraceRpg) []*entity.ChainTransfer {

	////
	filtertrace := []*util.TraceRpg{}
	for i, trace := range traces {
		if trace.Type == "call" {
			if trace.Value != "0" {
				trace.TraceTag = trace.TraceTag + "_" + gconv.String(i)
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
			if trace.Action.Value.String() != "0x0" && trace.Action.Input == "0x" && len(trace.TraceAddress) > 0 {
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
			Value:     trace.Action.Value.ToInt().String(),
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

// /////
func filteCalls(txIdx int, calls []*util.DebugTraceCalls) []*util.DebugTraceCalls {
	filtecalls := []*util.DebugTraceCalls{}
	for i, call := range calls {
		if len(call.Calls) > 0 {
			subfiltecalls := filteCalls(txIdx, call.Calls)
			if len(subfiltecalls) > 0 {
				for _, subcall := range subfiltecalls {
					subcall.TraceAddress = append(subcall.TraceAddress, i)
				}
				filtecalls = append(filtecalls, subfiltecalls...)
			}
		} else {
			if call.Type == "call" && call.Value.String() != "0x0" {
				call.TxIdx = txIdx
				call.TraceAddress = append(call.TraceAddress, i)
				filtecalls = append(filtecalls, call)
			}
		}
	}
	return filtecalls
}
func filteCallTrace(traces []*util.DebugTraceResult) []*util.DebugTraceCalls {
	filtetrace := []*util.DebugTraceCalls{}
	for i, trace := range traces {
		/////for subcall
		calls := filteCalls(i, trace.Result.Calls)
		///
		filtetrace = append(filtetrace, calls...)
	}
	return filtetrace
}
func ProcessInTxns_mantle(ctx context.Context, chainId int64, block *types.Block, traces []*util.DebugTraceResult) []*entity.ChainTransfer {
	////
	filtertrace := filteCallTrace(traces)

	//// fill transfer
	transfers := []*entity.ChainTransfer{}
	txs := block.Transactions()
	for _, trace := range filtertrace {
		tx := txs[trace.TxIdx]
		if tx == nil {
			g.Log().Fatal(ctx, "tx is nil:", "blockNumber:", block.Number())
			continue
		}
		/////
		transfer := &entity.ChainTransfer{
			ChainId:   chainId,
			Height:    block.Number().Int64(),
			BlockHash: block.Hash().String(),
			Ts:        int64(block.Time()),
			TxHash:    tx.Hash().String(),
			TxIdx:     trace.TxIdx,
			From:      trace.From.String(),
			To:        trace.To.String(),
			Contract:  "",
			Value:     trace.Value.ToInt().String(),
			Gas:       trace.Gas,
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
