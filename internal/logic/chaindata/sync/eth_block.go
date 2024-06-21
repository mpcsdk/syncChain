package block

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"syncChain/internal/conf"
	"syncChain/internal/logic/chaindata/sync/transfer"
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

func (s *EthModule) syncBlock() {
	defer func() {
		s.blockTimer.Reset(blockWait)
	}()

	client := s.getClient()
	if nil == client {
		s.logger.Errorf(s.ctx, "fail to get client")
		return
	}

	header := s.getHeader(client)
	if nil == header {
		g.Log().Error(s.ctx, "fail to get header")
		return
	}

	topHeight := header.Number.Int64()
	s.headerBlock = topHeight
	if s.lastBlock == 0 {
		s.lastBlock = topHeight
	}
	s.logger.Debugf(s.ctx, "chainId:%d, get header. height: %d, hash: %s", s.chainId, topHeight, header.Hash().String())

	if s.lastBlock >= s.headerBlock {
		s.logger.Infof(s.ctx, "no need to syncBlock, remote: %d, local: %d", topHeight, s.lastBlock)
		return
	}
	////
	////
	for {
		if topHeight > s.lastBlock {
			////batch proccess 10block
			startNumber := s.lastBlock + 1
			endNumber := s.lastBlock + 10
			if endNumber > topHeight {
				endNumber = topHeight
			}
			////
			wg := sync.WaitGroup{}
			lock := sync.Mutex{}
			//////
			txsmap := map[int64][]*entity.ChainTransfer{}
			errmap := map[int64]error{}
			///
			g.Log().Infof(s.ctx, "%d:syncBlock, startNumber: %d, endNumber: %d", s.chainId, startNumber, endNumber)
			for i := startNumber; i <= endNumber; i++ {
				wg.Add(1)
				go func(blockNumber int64) {
					defer wg.Done()
					txs, err := s.processBlock(s.ctx, blockNumber, client)
					if err != nil {
						lock.Lock()
						errmap[blockNumber] = err
						lock.Unlock()
					} else {
						lock.Lock()
						txsmap[blockNumber] = txs
						lock.Unlock()
					}
				}(i)
			}
			wg.Wait()
			//////
			if len(errmap) > 0 {
				return
			}
			////
			for _, v := range txsmap {
				s.transferCh <- v
			}
			s.lastBlock = endNumber
		} else {
			return
		}
	}

}

func (s *EthModule) processBlock(ctx context.Context, blockNumber int64, client *util.Client) ([]*entity.ChainTransfer, error) {
	transfers := []*entity.ChainTransfer{}
	block, _, txFroms, txHashes := s.getBlock(blockNumber, client)
	if nil == block {
		s.logger.Error(ctx, "fail to get block:", s.chainId, blockNumber)
		return nil, errors.New("fail to get block")
	}
	s.logger.Debugf(ctx, "chainId:%d , start getting blocks:%d", s.chainId, blockNumber)
	////get  transfers
	///process external
	for index, tx := range block.Transactions() {
		value := tx.Value()
		if tx == nil || tx.To() == nil || 0 == value.Sign() {
			continue
		}
		tx := transfer.ProcessTx(ctx, s.chainId, block, tx, txFroms, txHashes, index)
		if tx != nil {
			transfers = append(transfers, tx)
		}
	}
	s.logger.Debugf(ctx, "getTransaction,chainId:%d , number:%d, tx:%v", s.chainId, blockNumber, len(transfers))
	///internal
	if s.contracts.Len() != 0 {
		logs, err := s.getLogs(blockNumber, client)
		if err != nil {
			return nil, err
		}
		if len(logs) > 0 {
			txs := s.processEvent(int64(block.Time()), logs)
			transfers = append(transfers, txs...)
		}
		s.logger.Debugf(ctx, "getLogs,chainId:%d , number:%d, log:%d", s.chainId, blockNumber, len(logs))
	}
	/////
	////
	/////fake reciept
	for _, t := range transfers {
		t.Status = 1
		if t.ChainId == 0 {
			g.Log().Warning(ctx, t)
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

	return transfers, nil
}

func (s *EthModule) persistenceTransfer(txs []*entity.ChainTransfer) {
	/////
	g.Log().Debug(s.ctx, "persistenceTransfer:", s.chainId, s.lastBlock, txs)
	if len(txs) > 0 {
		// send latestTx
		service.EvnetSender().SendEvnetBatch_Latest(s.ctx, txs)
		////waiting for persistence
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
