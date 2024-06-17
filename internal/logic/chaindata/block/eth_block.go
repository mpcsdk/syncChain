package block

import (
	"context"
	"math/big"
	"syncChain/internal/conf"
	"syncChain/internal/logic/chaindata/block/transfer"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"
	"syncChain/internal/service"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

//	func skipToAddr(chainId int64, toaddr string) bool {
//		if addrs, ok := conf.Config.SkipToAddrChain[chainId]; ok {
//			if _, ok := addrs[toaddr]; ok {
//				return true
//			}
//		}
//		return false
//	}
func skipToAddr(chainId int64, toaddr string) bool {
	if addrs, ok := conf.Config.SkipToAddrChain[chainId]; ok {
		if _, ok := addrs[toaddr]; ok {
			return true
		}
	}
	return false
}

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
	s.headerBlock = topHeight
	if s.lastBlock == 0 {
		s.lastBlock = topHeight
	}
	s.logger.Debugf(s.ctx, "chainId:%d, get header. height: %d, hash: %s", s.chainId, topHeight, header.Hash().String())

	if s.lastBlock >= s.headerBlock {
		s.logger.Infof(s.ctx, "no need to processBlock, remote: %d, local: %d", topHeight, s.lastBlock)
		return
	}
	////
	////
	transfers := []*entity.ChainTransfer{}
	for i := s.lastBlock + 1; i < topHeight; i++ {
		block, _, txFroms, txHashes := s.getBlock(i, client)
		if nil == block {
			s.logger.Error(s.ctx, "fail to get block:", s.chainId)
			return
		}
		s.logger.Debugf(s.ctx, "chainId:%d , start getting blocks:%d:%d", s.chainId, i, block.NumberU64())
		////get  transfers
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
		/////
		////
		/////fake reciept
		for _, t := range transfers {
			t.Status = 1
			if t.ChainId == 0 {
				g.Log().Warning(s.ctx, t)
			}
		}

		///filter transfer of to in toaddrlist
		filtertransfer := []*entity.ChainTransfer{}
		for _, tx := range transfers {
			if tx.Kind == "erc20" || tx.Kind == "external" {
				if skipToAddr(s.chainId, tx.To) {
					continue
				}
			}
			filtertransfer = append(filtertransfer, tx)
		}
		transfers = filtertransfer
		if len(transfers) > 0 {
			// send latestTx
			service.EvnetSender().SendEvnetBatch_Latest(s.ctx, transfers)
			////waiting for persistence
			s.transferCh <- transfers
		}
		s.lastBlock = i
	}
}

func (s *EthModule) persistenceTransfer(txs []*entity.ChainTransfer) {
	/////
	if len(txs) > 0 {
		i := txs[0].Height
		s.blockTransfers[i] = txs
		s.logger.Debugf(s.ctx, "persistenceTransfer cached,chainId:%d , number:%d, log:%d", s.chainId, i, len(txs))
	}
	for i, txs := range s.blockTransfers {
		// /when last == topHeight - 12 insert into db
		if i > s.lastBlock-12 {
			err := service.DB().InsertTransferBatch(s.ctx, s.chainId, txs)
			if err != nil {
				if isDuplicateKeyErr(err) {
					s.logger.Warning(s.ctx, "fail to persistenceTransfer.  err:", err)
					err = service.DB().DelChainBlock(s.ctx, s.chainId, i)
					if err != nil {
						s.logger.Fatal(s.ctx, "fail to DelChainBlock. err:", err, txs)
						return
					}
					err = service.DB().InsertTransferBatch(s.ctx, s.chainId, txs)
				}
				if err != nil {
					s.logger.Fatal(s.ctx, "fail to persistenceTransfer. err: ", err, txs)
					return
				}
			}
			////send event
			service.EvnetSender().SendEvnetBatch(s.ctx, txs)
			s.updateHeight(i)
			delete(s.blockTransfers, i)
		}
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
