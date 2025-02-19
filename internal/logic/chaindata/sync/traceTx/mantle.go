package tracetx

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type MantleTrace struct {
	TraceSyncer
}

func newMantleTracer(ctx context.Context, chainId int64, url string, ctxTimeOut time.Duration) *MantleTrace {
	cli, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}
	return &MantleTrace{
		TraceSyncer: TraceSyncer{
			cli:        cli,
			ctx:        ctx,
			ctxTimeOut: ctxTimeOut,
			chainId:    chainId,
		},
	}
}

type DebugTraceResult struct {
	Result *DebugTraceCalls `json:"result"`
}
type DebugTraceCalls struct {
	From    common.Address     `json:"from"`
	Gas     string             `json:"gas"`
	GasUsed string             `json:"gasUsed"`
	To      common.Address     `json:"to"`
	Input   string             `json:"input"`
	Calls   []*DebugTraceCalls `json:"calls"`
	Value   *hexutil.Big       `json:"value"`
	Type    string             `json:"type"`
	////
	TxHash       common.Hash
	TxIdx        int
	TraceAddress []int
	///
	tag string
}

func (s *DebugTraceCalls) Tag() string {
	if s.tag == "" {
		s.tag = s.Type + "_" + gstr.JoinAny(s.TraceAddress, "_")
	}
	return s.tag
}
func (s *MantleTrace) GetTraceTransfer(ctx context.Context, block *ethtypes.Block) ([]*entity.SyncchainChainTransfer, error) {
	g.Log().Debug(ctx, "getDebug_TraceBlock:", block.Number())

	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeOut)
	defer cancel()
	//
	traces, err := s.debug_TraceBlock(ctx, block.Number())
	if err != nil {
		return nil, errors.New(fmt.Sprintln("getDebug_TraceBlock:", block.Number(), err))
	}
	t := s.processInTxns_mantle(ctx, block, traces)
	return t, nil
}

var callTracer = "callTracer"

func (s *MantleTrace) debug_TraceBlock(ctx context.Context, number *big.Int) ([]*DebugTraceResult, error) {
	var data []*DebugTraceResult
	err := s.cli.Client().CallContext(ctx, &data, "debug_traceBlockByNumber", toBlockNumArg(number), &tracers.TraceConfig{
		Tracer: &callTracer,
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data, err
}
func filteCalls(txIdx int, calls []*DebugTraceCalls) []*DebugTraceCalls {
	filtecalls := []*DebugTraceCalls{}
	for i, call := range calls {
		if len(call.Calls) > 0 {
			subfiltecalls := filteCalls(txIdx, call.Calls)
			if len(subfiltecalls) > 0 {
				for _, subcall := range subfiltecalls {
					subcall.TraceAddress = append(subcall.TraceAddress, i)
				}
				filtecalls = append(filtecalls, subfiltecalls...)
			}
		}
		if call.Type == "CALL" && call.Value.String() != "0x0" {
			call.TxIdx = txIdx
			call.TraceAddress = append(call.TraceAddress, i)
			filtecalls = append(filtecalls, call)
		}
	}
	return filtecalls
}

func filteMantleTrace(traces []*DebugTraceResult) []*DebugTraceCalls {
	filtetrace := []*DebugTraceCalls{}
	for i, trace := range traces {
		/////for subcall
		calls := filteCalls(i, trace.Result.Calls)
		///
		filtetrace = append(filtetrace, calls...)
	}
	return filtetrace
}

func (s *MantleTrace) processInTxns_mantle(ctx context.Context, block *ethtypes.Block, traces []*DebugTraceResult) []*entity.SyncchainChainTransfer {
	////
	filtertrace := filteMantleTrace(traces)

	//// fill transfer
	transfers := []*entity.SyncchainChainTransfer{}
	txs := block.Transactions()
	for _, trace := range filtertrace {
		tx := txs[trace.TxIdx]
		if tx == nil {
			g.Log().Fatal(ctx, "tx is nil:", "blockNumber:", block.Number())
			continue
		}
		/////
		transfer := &entity.SyncchainChainTransfer{
			ChainId:   s.chainId,
			Height:    block.Number().Int64(),
			BlockHash: block.Hash().String(),
			Ts:        int64(block.Time()),
			TxHash:    tx.Hash().String(),
			TxIdx:     trace.TxIdx,
			From:      trace.From.String(),
			To:        trace.To.String(),
			Contract:  common.Address{}.String(),
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
