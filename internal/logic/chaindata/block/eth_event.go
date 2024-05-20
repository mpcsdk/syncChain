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
)

func (s *EthModule) processEvent(txHash common.Hash, ts int64, client *util.Client) {
	receipt := s.getReceipt(txHash, client)
	if nil == receipt {
		return
	}

	for _, log := range receipt.Logs {
		topic := log.Topics[0].String()
		s.logger.Debugf(s.ctx, "chainId:%d,block:%d, tx: %s, get transfer data: %v", s.chainId, log.BlockNumber, log.TxHash.String(), log.Data)

		switch topic {
		case transferTopic:
			if len(log.Topics) == 3 {
				err := event.Process20(s.ctx, s.chainId, ts, log, int64(receipt.Status))
				if err != nil {
					s.logger.Warningf(s.ctx, "fail to unpack data.  err: %s", err)
					continue
				}
			} else if len(log.Topics) == 4 {
				err := event.Process721(s.ctx, s.chainId, ts, log, int64(receipt.Status))
				if err != nil {
					s.logger.Warningf(s.ctx, "fail to unpack data.  err: %s", err)
					continue
				}
			} else {
				s.logger.Notice(s.ctx, "unknown transfer topic: ", log)
			}
		case signalTopic:
			err := event.Process1155Signal(s.ctx, s.chainId, ts, log, int64(receipt.Status))
			if err != nil {
				s.logger.Warningf(s.ctx, "fail to unpack data.  err: %s", err)
				continue
			}
		case mulTopic:
			err := event.Process1155Batch(s.ctx, s.chainId, ts, log, int64(receipt.Status))
			if err != nil {
				s.logger.Warningf(s.ctx, "fail to unpack data.  err: %s", err)
				continue
			}
		default:
			s.logger.Debug(s.ctx, "unknown event topic:", log)
		}
	}

}

func (s *EthModule) getReceipt(txHash common.Hash, client *util.Client) *types.Receipt {
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
			s.logger.Errorf(s.ctx, "fail to get logs, err: %s, close client and reconnect", err)
			s.closeClient()
			return nil
		}
		return receipt
	case <-ctx.Done():
		s.logger.Errorf(s.ctx, "fail to get logs, err: timeout, close client and reconnect")
		s.closeClient()
		return nil
	}
}

func (s *EthModule) getLogs(i int64, client *util.Client) []types.Log {
	var (
		logs []types.Log
		err  error
	)
	ch := make(chan byte, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		var query ethereum.FilterQuery
		query.FromBlock = big.NewInt(i)
		query.ToBlock = big.NewInt(i)
		query.Addresses = s.contracts.Addresses()
		logs, err = client.FilterLogs(ctx, query)

		ch <- 0
	}()

	select {
	case <-ch:
		if err != nil {
			s.logger.Errorf(s.ctx, "fail to get logs, err: %s, close client and reconnect", err)
			s.closeClient()
			return nil
		}

		// success, but no result
		if nil == logs {
			logs = []types.Log{}
		}
		return logs
	case <-ctx.Done():
		s.logger.Errorf(s.ctx, "fail to get logs, err: timeout, close client and reconnect")
		s.closeClient()
		return nil
	}
}
