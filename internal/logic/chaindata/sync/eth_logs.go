package syncBlock

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"syncChain/internal/logic/chaindata/sync/transfer"
	"syncChain/internal/logic/chaindata/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/lib/pq"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	// "github.com/lib/pq"
)

func isDuplicateKeyErr(err error) bool {
	gerr := err.(*gerror.Error)
	if cerr, ok := gerr.Cause().(*pq.Error); ok {
		if cerr.Code == "23505" {
			return true
		}
	}
	return false
}

func (s *EthModule) processEvent(ts int64, logs []types.Log) []*entity.SyncchainChainTransfer {
	txs := []*entity.SyncchainChainTransfer{}
	for _, log := range logs {

		topic := log.Topics[0].String()
		g.Log().Debug(s.ctx, "processEvent chainId:", s.chainId, "block:", log.BlockNumber, "tx:", log.TxHash.String(), "topic:", topic)

		switch topic {
		case transferTopic:
			if len(log.Topics) == 3 {
				tx := transfer.Process20(s.ctx, s.chainId, ts, &log)
				if tx == nil {
					g.Log().Error(s.ctx, "fail to Process20.  err:", log)
				} else {
					txs = append(txs, tx)
				}
			} else if len(log.Topics) == 4 {
				tx := transfer.Process721(s.ctx, s.chainId, ts, &log)
				if tx == nil {
					g.Log().Error(s.ctx, "fail to Process721.  err:", log)
				} else {
					txs = append(txs, tx)
				}
			} else {
				g.Log().Notice(s.ctx, "unknown transfer topic: ", log)
			}
		case signalTopic:
			tx := transfer.Process1155Signal(s.ctx, s.chainId, ts, &log)
			if tx == nil {
				g.Log().Error(s.ctx, "fail to Process1155Signal.  err:", log)
			} else {
				txs = append(txs, tx)
			}
		case mulTopic:
			tx := transfer.Process1155Batch(s.ctx, s.chainId, ts, &log)
			if tx == nil {
				g.Log().Error(s.ctx, "fail to Process1155Signal.  err:", log)
			}
			txs = append(txs, tx...)
		default:
			g.Log().Warning(s.ctx, "unknown event topic:", s.chainId, log)
		}
	}
	return txs

}

var topic [][]common.Hash = [][]common.Hash{
	{
		common.HexToHash(transferTopic),
		common.HexToHash(signalTopic),
		common.HexToHash(mulTopic),
	},
}

func (s *EthModule) getLogs(i int64) ([]types.Log, error) {
	g.Log().Debug(s.ctx, "getLogs:", s.chainId, i)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	logs, err := s.cli.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(i),
		ToBlock:   big.NewInt(i),
		Addresses: s.syncContracts,
		Topics:    topic,
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintln("getLogs:", i, err))
	}
	return logs, nil
}
