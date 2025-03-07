package tracetx

import (
	"context"
	"fmt"
	"math/big"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type ITraceSyncer interface {
	GetTraceTransfer(ctx context.Context, block *types.Block) ([]*entity.SyncchainChainTransfer, error)
}

type TraceSyncer struct {
	cli        *util.Client
	ctx        context.Context
	ctxTimeOut time.Duration
	chainId    int64
}

func NewTraceSyncer(ctx context.Context, chainId int64, url string, ctxTimeOut time.Duration) ITraceSyncer {
	cli, err := util.Dial(url)
	if err != nil {
		panic(err)
	}
	cid, err := cli.ChainID(ctx)
	if err != nil {
		panic(err)
	}
	if cid.Int64() != chainId {
		panic("chainId not match")
	}

	///
	switch chainId {
	// case 9527, 2025:
	// 	return newRpgTracer(ctx, chainId)
	case 1, 11155111, 97, 56, 5000, 5003:
		return newMantleTracer(ctx, chainId, url, ctxTimeOut)
	// case 5000, 5003:
	// 	return newMantleTracer(ctx, chainId, url, ctxTimeOut)
	default:
		return &Empty{}
	}
	return &Empty{}
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}
