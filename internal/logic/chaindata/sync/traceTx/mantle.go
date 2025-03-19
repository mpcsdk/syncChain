package tracetx

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type MantleTrace struct {
	TraceSyncer
}

func newMantleTracer(ctx context.Context, chainId int64, url string, ctxTimeOut time.Duration) *MantleTrace {
	cli, err := util.Dial(url)
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
	TxHash common.Hash      `json:"txHash"`
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
	Error   string             `json:"error"`
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
func (s *MantleTrace) GetTraceTransfer(ctx context.Context, block *types.Block) ([]*entity.SyncchainChainTransfer, error) {
	g.Log().Debug(ctx, "getDebug_TraceBlock:", block.Number())

	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeOut*5)
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
func filteCalls(calls []*DebugTraceCalls) []*DebugTraceCalls {
	filtecalls := []*DebugTraceCalls{}
	for i, call := range calls {
		if len(call.Calls) > 0 {
			subfiltecalls := filteCalls(call.Calls)
			if len(subfiltecalls) > 0 {
				for _, subcall := range subfiltecalls {
					subcall.TraceAddress = append(subcall.TraceAddress, i)
				}
				filtecalls = append(filtecalls, subfiltecalls...)
			}
		}
		if call.Type == "CALL" && call.Value.String() != "0x0" && call.Input == "0x" {
			call.TraceAddress = append(call.TraceAddress, i)
			filtecalls = append(filtecalls, call)
		}
	}
	return filtecalls
}

func filteMantleTrace(traces []*DebugTraceResult, txs []*types.Transaction) []*DebugTraceCalls {
	filtetrace := []*DebugTraceCalls{}
	for pos, trace := range traces {
		if trace.Result.Error != "" {
			continue
		}
		/////for subcall
		// txIdx := 0
		// pos := -1
		// if trace.TxHash == (common.Hash{}) {
		// 	pos = findTransactionPosByFromTo(trace.Result.From.Hex(), trace.Result.To.Hex(), txs)
		// } else {
		// 	pos = findTransactionPos(trace.TxHash, txs)
		// }

		if pos > len(txs) {
			g.Log().Fatal(context.Background(), "trace not find tx:\n", trace, "\ntxs:", txs)
			continue
		}
		/////
		calls := filteCalls(trace.Result.Calls)
		for _, call := range calls {
			call.TxHash = txs[pos].Hash()
			call.TxIdx = pos
		}
		/////
		///
		filtetrace = append(filtetrace, calls...)
	}
	return filtetrace
}
func findTransactionPos(hash common.Hash, txs []*types.Transaction) int {
	for i, tx := range txs {
		if tx.Hash().String() == hash.String() {
			return i
		}
	}
	return -1
}
func findTransactionPosByFromTo(from, to string, txs []*types.Transaction) int {
	for i, tx := range txs {
		if tx.Sender().String() == from && (tx.To() == nil || tx.To().String() == to) {
			return i
		}
	}
	return -1
}
func (s *MantleTrace) processInTxns_mantle(ctx context.Context, block *types.Block, traces []*DebugTraceResult) []*entity.SyncchainChainTransfer {
	////
	txs := block.Transactions()
	filtertrace := filteMantleTrace(traces, txs)

	//// fill transfer
	transfers := []*entity.SyncchainChainTransfer{}
	for _, trace := range filtertrace {
		tx := txs[trace.TxIdx]
		transfer := &entity.SyncchainChainTransfer{
			ChainId:   s.chainId,
			Height:    block.Number().Int64(),
			BlockHash: block.Hash().String(),
			Ts:        int64(block.Time()),
			TxHash:    trace.TxHash.String(),
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
		if len(transfer.TraceTag) > 20 {
			g.Log().Debug(ctx, "toolong trace:", transfer)
			continue
		}
		transfers = append(transfers, transfer)
	}

	return transfers
}
