package tracetx

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

var rpgtraceurl = "https://mainnet.rangersprotocol.com/api"
var rpgtraceurl_testnet = "https://robin-api.rangersprotocol.com"

type RpgTrace struct {
	TraceSyncer
}

type RpgTraceResult struct {
	Data []*TraceRpg `json:"data"`
}
type TraceRpg struct {
	BlockHeight  int64          `json:"blockHeight"`
	BlockHash    common.Hash    `json:"blockHash"`
	Depth        int            `json:"depth"`
	GasLimit     string         `json:"gasLimit"`
	ParentTxHash common.Hash    `json:"parentTxHash"`
	TxIndex      int            `json:"txIndex"`
	Source       common.Address `json:"source"`
	Target       common.Address `json:"target"`
	Time         time.Time      `json:"time"`
	Type         string         `json:"type"`
	Value        string         `json:"value"`
	TraceTag     string         `json:"traceTag"`
}

func newRpgTracer(ctx context.Context, chainId int64) *RpgTrace {
	if chainId == 9527 {
		cli, err := ethclient.Dial(rpgtraceurl_testnet)
		if err != nil {
			panic(err)
		}
		return &RpgTrace{
			TraceSyncer: TraceSyncer{
				ctx:        ctx,
				ctxTimeOut: time.Second * 10,
				cli:        cli,
				chainId:    chainId,
			},
		}
	} else {
		cli, err := ethclient.Dial(rpgtraceurl)
		if err != nil {
			panic(err)
		}
		return &RpgTrace{
			TraceSyncer: TraceSyncer{
				ctx:        ctx,
				ctxTimeOut: time.Second * 10,
				cli:        cli,
				chainId:    chainId,
			},
		}

	}

}
func (s *RpgTrace) GetTraceTransfer(ctx context.Context, block *ethtypes.Block) ([]*entity.SyncchainChainTransfer, error) {
	g.Log().Debug(ctx, "getTraceBlock_rpg:", block.Number())

	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeOut)
	defer cancel()
	//
	traces, err := s.traceBlock_rpg(ctx, block.Number())
	if err != nil {
		return nil, errors.New(fmt.Sprintln("getTraceBlock_rpg:", block.Number(), err))
	}
	t := s.ProcessInTxnsRpg(ctx, block, traces)
	return t, nil
}
func (s *RpgTrace) traceBlock_rpg(ctx context.Context, number *big.Int) ([]*TraceRpg, error) {
	var data *RpgTraceResult
	err := s.cli.Client().CallContext(ctx, &data, "Rocket_getInternalTxByBlock", number)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data.Data, err
}
func (s *RpgTrace) ProcessInTxnsRpg(ctx context.Context, block *ethtypes.Block, traces []*TraceRpg) []*entity.SyncchainChainTransfer {
	////
	filtertrace := []*TraceRpg{}
	for i, trace := range traces {
		if trace.Type == "call" {
			if trace.Value != "0" {
				trace.TraceTag = trace.TraceTag + "_" + gconv.String(i)
				filtertrace = append(filtertrace, trace)
			}
		}
	}

	//// fill transfer
	transfers := []*entity.SyncchainChainTransfer{}
	for _, trace := range filtertrace {
		tx := block.Transaction(trace.ParentTxHash)
		if tx == nil {
			g.Log().Warning(ctx, "tx is nil")
			continue
		}
		/////
		transfer := &entity.SyncchainChainTransfer{
			ChainId:   s.chainId,
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
