package block

import (
	"context"
	"math/big"
	"syncChain/internal/logic/chaindata/block/event"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"
	"time"

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
func (s *EthModule) processEvent(ts int64, logs []types.Log) []*entity.ChainTransfer {

	txs := []*entity.ChainTransfer{}
	for _, log := range logs {
		topic := log.Topics[0].String()
		s.logger.Debug(s.ctx, "processEvent chainId:", s.chainId, "block:", log.BlockNumber, "tx:", log.TxHash.String(), "topic:", topic)

		switch topic {
		case transferTopic:
			if len(log.Topics) == 3 {
				tx := event.Process20(s.ctx, s.chainId, ts, &log)
				if tx == nil {
					g.Log().Error(s.ctx, "fail to Process20.  err:", log)
				}
				txs = append(txs, tx)
			} else if len(log.Topics) == 4 {
				tx := event.Process721(s.ctx, s.chainId, ts, &log)
				if tx == nil {
					g.Log().Error(s.ctx, "fail to Process721.  err:", log)
				}
				txs = append(txs, tx)
			} else {
				s.logger.Notice(s.ctx, "unknown transfer topic: ", log)
			}
		case signalTopic:
			tx := event.Process1155Signal(s.ctx, s.chainId, ts, &log)
			if tx == nil {
				g.Log().Error(s.ctx, "fail to Process1155Signal.  err:", log)
			}
			txs = append(txs, tx)
		case mulTopic:
			tx := event.Process1155Batch(s.ctx, s.chainId, ts, &log)
			if tx == nil {
				g.Log().Error(s.ctx, "fail to Process1155Signal.  err:", log)
			}
			txs = append(txs, tx...)
		default:
			s.logger.Debug(s.ctx, "unknown event topic:", log)
		}
	}
	return txs

}

func (s *EthModule) getReceipt(txHash common.Hash, client *util.Client) *types.Receipt {
	g.Log().Debug(s.ctx, "eth_getReceipt:", s.chainId, txHash)
	var (
		err     error
		receipt *types.Receipt
	)
	ch := make(chan byte, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		// var query ethereum.FilterQuery
		// query.FromBlock = big.NewInt(i)
		// query.ToBlock = big.NewInt(i)
		// query.Addresses = s.contracts.Addresses()
		receipt, err = client.TransactionReceipt(ctx, txHash)

		ch <- 0
	}()

	select {
	case <-ch:
		if err != nil {
			if err.Error() == "not found" {
				return &types.Receipt{
					Status: types.ReceiptStatusFailed,
				}
			}
			s.logger.Error(s.ctx, "fail to TransactionReceipt:", txHash, "err:", err)
			s.closeClient()
			return nil
		}
		return receipt
	case <-ctx.Done():
		s.logger.Errorf(s.ctx, "fail to get TransactionReceipt, err: timeout, close client and reconnect")
		s.closeClient()
		return nil
	}
}

func (s *EthModule) getLogs(i int64, client *util.Client) []types.Log {
	g.Log().Debug(s.ctx, "eth_getLogs:", s.chainId, i)
	var (
		logs []types.Log
		err  error
	)
	ch := make(chan byte, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query ethereum.FilterQuery
	go func() {
		query.FromBlock = big.NewInt(i)
		query.ToBlock = big.NewInt(i)
		query.Addresses = s.contracts.Addresses()
		logs, err = client.FilterLogs(ctx, query)

		ch <- 0
	}()

	select {
	case <-ch:
		if err != nil {
			s.logger.Error(s.ctx, "fail to get logs,chain:", s.chainId, "query:", query, "err:", err)
			s.closeClient()
			return nil
		}

		// success, but no result
		if nil == logs {
			logs = []types.Log{}
		}
		return logs
	case <-ctx.Done():
		s.logger.Error(s.ctx, "fail to get logs, err: timeout, close client and reconnect:", s.chainId)
		s.closeClient()
		return nil
	}
}
