package block

import (
	"context"
	"math/big"
	"syncChain/internal/logic/chaindata/block/transfer"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"
	"syncChain/internal/service"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

// func skipToAddr(chainId int64, toaddr string) bool {
// 	if addrs, ok := conf.Config.SkipToAddrChain[chainId]; ok {
// 		if _, ok := addrs[toaddr]; ok {
// 			return true
// 		}
// 	}
// 	return false
// }

func (s *EthModule) processBlock() {
	s.lock.Lock()
	defer func() {
		s.blockTimer.Reset(blockWait)
		s.lock.Unlock()
	}()

	client := s.getClient()
	if nil == client {
		s.logger.Errorf(s.ctx, "fail to get client")
		return
	}

	header := s.getHeader(client)
	if nil == header {
		return
	}

	topHeight := header.Number.Int64()
	if s.lastBlock == 0 {
		s.lastBlock = topHeight
	}
	s.logger.Debugf(s.ctx, "chainId:%d, get header. height: %d, hash: %s", s.chainId, topHeight, header.Hash().String())

	last := topHeight - 12
	if last == s.lastBlockFromClient {
		s.count++
		if s.count == 6 {
			s.logger.Warningf(s.ctx, "get max retry count, close client and reconnect")
			s.closeClient()
			return
		}
	} else {
		s.count = 0
		s.lastBlockFromClient = last
	}

	if last <= s.lastBlock {
		s.logger.Infof(s.ctx, "no need to processBlock, remote: %d, local: %d", last, s.lastBlock)
		return
	}

	for i := s.lastBlock + 1; i < last; i++ {
		block, _, txFroms, txHashes := s.getBlock(i, client)
		if nil == block {
			s.logger.Error(s.ctx, "fail to get block:", s.chainId)
			return
		}
		s.logger.Debugf(s.ctx, "chainId:%d , start getting blocks:%d:%d", s.chainId, i, block.NumberU64())
		////get  transfers
		transfers := []*entity.ChainTransfer{}

		///process external
		for index, tx := range block.Transactions() {
			value := tx.Value()
			if tx == nil || tx.To() == nil || 0 == value.Sign() {
				continue
			}
			tx := transfer.ProcessTx(s.ctx, s.chainId, block, tx, txFroms, txHashes, index)
			if tx != nil {
				transfers = append(transfers, tx)
			}
		}
		s.logger.Debugf(s.ctx, "getTransaction,chainId:%d , number:%d, tx:%v", s.chainId, i, len(transfers))
		///internal
		if s.contracts.Len() != 0 {
			logs := s.getLogs(i, client)
			if len(logs) > 0 {
				txs := s.processEvent(int64(block.Time()), logs)
				transfers = append(transfers, txs...)
			}
			s.logger.Debugf(s.ctx, "getLogs,chainId:%d , number:%d, log:%d", s.chainId, i, len(logs))
		}
		///get receipt
		// for i, tx := range transfers {
		// 	receipt := s.getReceipt(common.HexToHash(tx.TxHash), client)
		// 	if nil == receipt {
		// 		receipt = &types.Receipt{
		// 			Status: types.ReceiptStatusFailed,
		// 		}
		// 	} else {
		// 		if receipt.TxHash.Hex() != tx.TxHash {
		// 			receipt.Status = types.ReceiptStatusFailed
		// 		}
		// 	}
		// 	transfers[i].Status = int64(receipt.Status)
		// }
		/////fake reciept
		for _, t := range transfers {
			t.Status = 1
			if t.ChainId == 0 {
				g.Log().Warning(s.ctx, t)
			}
		}

		///insert transfers
		////
		if len(transfers) != 0 {
			err := service.DB().InsertTransferBatch(s.ctx, s.chainId, transfers)
			if err != nil {
				if isDuplicateKeyErr(err) {
					s.logger.Warning(s.ctx, "fail to InsertTransferBatch.  err:", err)
					return
				}
				s.logger.Fatal(s.ctx, "fail to InsertTransferBatch.  err: ", err)
				return
			}
		}
		s.logger.Debugf(s.ctx, "InsertTransfer,chainId:%d , number:%d, log:%d", s.chainId, i, len(transfers))

		////
		s.lastBlock = i
		s.updateHeight()
	}
}

func (s *EthModule) getBlock(i int64, client *util.Client) (*types.Block, *common.Hash, []*common.Address, []*common.Hash) {
	g.Log().Debug(s.ctx, "eth_getBlock:", s.chainId, i)
	var (
		block    *types.Block
		hash     *common.Hash
		txFroms  []*common.Address
		txHashes []*common.Hash
		err      error
	)
	ch := make(chan byte, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		block, hash, txFroms, txHashes, err = client.BlockByNumber(ctx, big.NewInt(i))
		ch <- 0
	}()

	select {
	case <-ch:
		if err != nil {
			s.logger.Error(s.ctx, "fail to get block:", s.chainId, "err:", err)
			s.closeClient()
			return nil, nil, nil, nil
		}

		return block, hash, txFroms, txHashes
	case <-ctx.Done():
		s.logger.Error(s.ctx, "fail to get block:", s.chainId, " timeout, close client and reconnect")
		s.closeClient()
		return nil, nil, nil, nil
	}
}
