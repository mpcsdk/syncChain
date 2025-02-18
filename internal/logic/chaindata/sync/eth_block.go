package syncBlock

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"syncChain/internal/conf"
	"syncChain/internal/logic/chaindata/sync/transfer"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"
	"syncChain/internal/service"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func (s *EthModule) isSkipToAddr(toaddr string) bool {
	if _, ok := s.skipToAddrs[toaddr]; ok {
		return true
	}
	return false
}
func (s *EthModule) isSkipFromAddr(addr string) bool {
	if _, ok := s.skipFromAddrs[addr]; ok {
		return true
	}
	return false
}

func (s *EthModule) syncBlock() {
	defer func() {
		s.blockTimer.Reset(s.blockWait)
	}()

	client := s.cli
	if nil == client {
		g.Log().Errorf(s.ctx, "fail to get client")
		return
	}

	nr, err := s.getBlockNumber(client)
	if err != nil {
		g.Log().Error(s.ctx, "fail to get header")
		return
	}

	latestBlock := nr
	topHeight := latestBlock - conf.Config.Syncing.WaitBlock
	g.Log().Infof(s.ctx, "chainId:%d, get header. latest: %d, topHeight: %d", s.chainId, latestBlock, topHeight)
	////
	//// syncbatchblock
	for {
		if topHeight > s.currentBlock {
			////batch proccess 10block
			startNumber := s.currentBlock + 1
			if startNumber >= topHeight {
				g.Log().Infof(s.ctx, "no need to syncBlock, remote: %d, local: %d", topHeight, s.currentBlock)
				return
			}

			endNumber := s.currentBlock + conf.Config.Syncing.BatchSyncTask
			if endNumber > topHeight {
				endNumber = topHeight
			}
			////
			wg := sync.WaitGroup{}
			lock := sync.Mutex{}
			//////
			txsmap := map[int64][]*entity.SyncchainChainTransfer{}
			errmap := map[int64]error{}
			///
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
				for k, v := range errmap {
					g.Log().Error(s.ctx, "batchSync err:", k, v)
				}
				return
			}
			////sortmap
			sortnuber := []int64{}
			for i, _ := range txsmap {
				sortnuber = append(sortnuber, i)
			}
			slices.Sort(sortnuber)
			////
			cnt := 0
			for _, v := range sortnuber {
				cnt = cnt + len(txsmap[v])
				s.transferCh <- txsmap[v]
			}
			g.Log().Infof(s.ctx, "%d:syncBlock, startNumber: %d, endNumber: %d, cnt:%d", s.chainId, startNumber, endNumber, cnt)
			s.currentBlock = endNumber
			// s.lastBlock = endNumber
		} else {
			return
		}
	}

}

var rpgAddr = common.HexToAddress("0x71d9CFd1b7AdB1E8eb4c193CE6FFbe19B4aeE0dB").String()

func (s *EthModule) processBlock(ctx context.Context, blockNumber int64, client *ethclient.Client) ([]*entity.SyncchainChainTransfer, error) {
	transfers := []*entity.SyncchainChainTransfer{}
	// block, _, txFroms, txHashes, err := s.getBlockByNumber(blockNumber, client)
	block, err := s.getBlockByNumber(blockNumber, client)
	if err != nil {
		return nil, err
	}
	if nil == block {
		g.Log().Error(ctx, "fail to get block:", s.chainId, blockNumber)
		return nil, errors.New(fmt.Sprintln("fail to get block:", blockNumber))
	}
	////get  transfers
	///process external
	for index, tx := range block.Transactions() {
		value := tx.Value()
		if tx == nil || tx.To() == nil || 0 == value.Sign() {
			continue
		}
		tx := transfer.ProcessBlock(ctx, block, tx, index)
		if tx != nil {
			transfers = append(transfers, tx)
		}
	}
	g.Log().Debugf(ctx, "getTransaction,chainId:%d , number:%d, tx:%v", s.chainId, blockNumber, len(transfers))
	/////notice: trace_block Internal Txns
	traceTransfer, err := s.tracer.GetTraceTransfer(ctx, block)
	if err != nil {
		g.Log().Error(ctx, "fail to getTraceBlock:", blockNumber, err)
		return nil, errors.New(fmt.Sprintln("fail GetTraceTransfer:", blockNumber))
	}
	transfers = append(transfers, traceTransfer...)
	///logs
	if len(s.syncContracts) != 0 {
		logs, err := s.getLogs(blockNumber)
		if err != nil {
			return nil, err
		}
		if len(logs) > 0 {
			txs := s.processEvent(int64(block.Time()), logs)
			transfers = append(transfers, txs...)
		}
		g.Log().Debugf(ctx, "getLogs,chainId:%d , number:%d, log:%d", s.chainId, blockNumber, len(logs))
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
	///and rpg contract to native
	filtertransfer := []*entity.SyncchainChainTransfer{}
	for _, tx := range transfers {
		if tx.Kind == "erc20" {
			if s.isSkipToAddr(tx.To) || s.isSkipFromAddr(tx.From) {
				continue
			}
			if tx.Contract == rpgAddr {
				tx.Contract = ""
				tx.Kind = "external"
			}
		} else if tx.Kind == "external" {
			if s.isSkipToAddr(tx.To) || s.isSkipFromAddr(tx.From) {
				continue
			}
		}
		filtertransfer = append(filtertransfer, tx)
	}
	transfers = filtertransfer

	return transfers, nil
}

func (s *EthModule) persistenceTransfer(txs []*entity.SyncchainChainTransfer) {
	/////
	g.Log().Debug(s.ctx, "persistenceTransfer:", s.chainId, s.currentBlock, txs)
	if len(txs) > 0 {
		// send latestTx
		service.EvnetSender().SendEvnetBatch_Latest(s.ctx, txs)
		////waiting for persistence
		i := txs[0].Height
		s.blockTransfers[i] = txs
		g.Log().Debugf(s.ctx, "persistenceTransfer cached,chainId:%d , number:%d, log:%d", s.chainId, i, len(txs))
	}
	for i, txs := range s.blockTransfers {
		// wait 12 block when block is confirmed
		if i > s.currentBlock-12 {
			err := service.DB().InsertTransfer_Transaction(s.ctx, s.chainId, txs)
			if err != nil {
				g.Log().Fatal(s.ctx, "InsertTransfer_Transaction:", err)
				// if isDuplicateKeyErr(err) {
				// 	g.Log().Warning(s.ctx, "fail to persistenceTransfer.  err:", err)
				// 	err = service.DB().DelChainBlock(s.ctx, s.chainId, i)
				// 	if err != nil {
				// 		g.Log().Fatal(s.ctx, "fail to DelChainBlock. err:", err, txs)
				// 		return
				// 	}
				// 	err = service.DB().InsertTransferBatch(s.ctx, s.chainId, txs)
				// }
				// if err != nil {
				// 	g.Log().Fatal(s.ctx, "fail to persistenceTransfer. err: ", err)
				// 	return
				// }
			}
			////send event
			service.EvnetSender().SendEvnetBatch(s.ctx, txs)
			s.updateHeight(i)
			delete(s.blockTransfers, i)
		}
	}

}
func (s *EthModule) getBlockNumber(client *ethclient.Client) (int64, error) {
	g.Log().Debug(s.ctx, "eth_blockNumber:", s.chainId)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	nr, err := client.BlockNumber(ctx)

	if err != nil {
		return 0, errors.New(fmt.Sprintln("eth_blockNumber:", err))
	}
	return int64(nr), nil
}

func (s *EthModule) getBlock(i int64, client *util.Client) (*types.Block, *common.Hash, []*common.Address, []*common.Hash, error) {
	g.Log().Debug(s.ctx, "eth_getBlock:", s.chainId, i)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	block, hash, txFroms, txHashes, err := client.BlockByNumber(ctx, big.NewInt(i))

	if err != nil {
		return nil, nil, nil, nil, errors.New(fmt.Sprintln("eth_getBlock:", i, err))
	}
	return block, hash, txFroms, txHashes, nil
}
func (s *EthModule) getBlockByNumber(i int64, client *ethclient.Client) (*ethtypes.Block, error) {
	g.Log().Debug(s.ctx, "eth_getBlock:", s.chainId, i)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	block, err := client.BlockByNumber(ctx, big.NewInt(i))

	if err != nil {
		return nil, errors.New(fmt.Sprintln("eth_getBlock:", i, err))
	}
	return block, nil
}

func (s *EthModule) getTraceBlock(i int64, client *util.Client) ([]*util.Trace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	traces, err := client.TraceBlock(ctx, big.NewInt(i))
	if err != nil {

		return nil, errors.New(fmt.Sprintln("getTraceBlock:", i, err))

	}
	return traces, nil
}
func (s *EthModule) getTraceBlock_rpg(i int64, client *util.Client) ([]*util.TraceRpg, error) {
	g.Log().Debug(s.ctx, "getTraceBlock_rpg:", s.chainId, i)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()
	//
	traces, err := client.TraceBlock_rpg(ctx, i)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("getTraceBlock_rpg:", i, err))
	}
	return traces, nil
}
func (s *EthModule) getDebug_TraceBlock(i int64, client *util.Client) ([]*util.DebugTraceResult, error) {
	g.Log().Debug(s.ctx, "getDebug_TraceBlock:", s.chainId, i)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()
	//
	traces, err := client.Debug_TraceBlock(ctx, big.NewInt(i))
	if err != nil {
		return nil, errors.New(fmt.Sprintln("getDebug_TraceBlock:", i, err))
	}
	return traces, nil
}
func (s *EthModule) getHeader(client *util.Client) (*types.Header, error) {
	g.Log().Debug(s.ctx, "eth_Header:", s.chainId)
	// var (
	// 	header *types.Header
	// 	err    error
	// )
	// ch := make(chan byte, 1)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	header, err := client.HeaderByNumber(ctx, nil)

	return header, err
}
