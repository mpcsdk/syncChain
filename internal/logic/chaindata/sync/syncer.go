package syncBlock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"syncChain/internal/conf"
	tracetx "syncChain/internal/logic/chaindata/sync/traceTx"
	"syncChain/internal/logic/chaindata/util"
	"syncChain/internal/service"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type EthModule struct {
	ctx context.Context
	// cli    *ethclient.Client
	cli    *util.Client
	tracer tracetx.ITraceSyncer
	///
	isRunning bool
	exitSig   bool
	// exit   chan bool
	// pause  chan bool
	// closed bool
	///
	rpcUrl string
	// name    string
	chainId int64

	////
	syncContracts []common.Address
	skipToAddrs   map[string]struct{}
	skipFromAddrs map[string]struct{}

	// client        *util.Client
	// last block from client
	// last block processed
	lastBlock    int64
	currentBlock int64
	// confirmedBlock int64
	// startBlock int64

	// blockTimer *time.Timer
	blockWait time.Duration
	// = time.Duration(conf.Config.Server.SyncInterval) * time.Second
	// clientTimer *time.Timer

	// abi           abi.ABI
	// event         abi.Event
	// transferTopic string
	// lock sync.Mutex
	//
	// logger     *glog.Logger
	// chaincfgdb *mpcdao.ChainCfg
	////
	transferCh     chan []*entity.SyncchainChainTransfer
	blockTransfers map[int64][]*entity.SyncchainChainTransfer
	///
}

// var rpgtraceurl = "https://mainnet.rangersprotocol.com/api"
// var rpgtraceurl_testnet = "https://robin-api.rangersprotocol.com"

func NewEthModule(ctx context.Context, chainId int64, currentBlock int64, rpcUrl string, syncContracts []common.Address, skipToAddrs []common.Address, skipFromAddrs []common.Address) *EthModule {
	cli, err := util.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}
	if currentBlock == 0 {
		nr, err := cli.BlockNumber(ctx)
		if err != nil {
			panic(err)
		}
		currentBlock = int64(nr)
	}
	///
	tracer := tracetx.NewTraceSyncer(
		ctx, chainId,
		rpcUrl,
		time.Duration(conf.Config.Syncing.TimeOut)*time.Second,
	)
	s := &EthModule{
		ctx:          ctx,
		chainId:      chainId,
		currentBlock: currentBlock,
		// startBlock:   currentBlock,
		// lastBlock: currentBlock,
		rpcUrl: rpcUrl,
		cli:    cli,
		tracer: tracer,
		// exit:           make(chan bool),
		// pause:          make(chan bool),
		// closed:         false,
		syncContracts: syncContracts,
		skipToAddrs: func() map[string]struct{} {
			addrs := map[string]struct{}{}
			for _, addr := range skipToAddrs {
				addrs[addr.String()] = struct{}{}
			}
			return addrs

		}(),
		skipFromAddrs: func() map[string]struct{} {
			addrs := map[string]struct{}{}
			for _, addr := range skipFromAddrs {
				addrs[addr.String()] = struct{}{}
			}
			return addrs

		}(),
		// chaincfgdb:     mpcdao.NewChainCfg(nil, 0),
		////
		transferCh:     make(chan []*entity.SyncchainChainTransfer, 100),
		blockTransfers: map[int64][]*entity.SyncchainChainTransfer{},
		///
		blockWait: time.Duration(conf.Config.Syncing.BlockInterval) * time.Second,
	}

	////
	// s.blockTimer = time.NewTimer(s.blockWait)
	// s.blockTimer.Stop()
	///
	// s.loop()
	//
	return s
}
func (s *EthModule) Exit() {
	s.exitSig = true
}
func (s *EthModule) IsRunning() bool {
	return s.isRunning
}

func (s *EthModule) Start() {
	s.isRunning = true
	s.exitSig = false
	client := s.cli
	if nil == client {
		g.Log().Fatal(s.ctx, "fail to get client")
		return
	}
	_, err := s.cli.BlockNumber(s.ctx)
	if err != nil {
		g.Log().Fatal(s.ctx, "fail to get header")
		return
	}
	go func() {
		for {
			if s.exitSig {
				s.isRunning = false
				return
			}
			nr, err := s.getBlockNumber(client)
			if err != nil {
				g.Log().Error(s.ctx, "fail to get header")
			} else {
				err = s.syncBlock(nr)
			}
			////
			if err != nil || s.currentBlock >= nr {
				time.Sleep(s.blockWait)
			}
			if nr%1000 == 0 {
				err := service.DB().TruncateTransfer(s.ctx, s.chainId, s.currentBlock-conf.Config.Syncing.PersistenceBlocks)
				g.Log().Info(s.ctx, "truncate transfer, block:", s.currentBlock-conf.Config.Syncing.PersistenceBlocks)
				if err != nil {
					g.Log().Error(s.ctx, "fail to truncate transfer, err:", err)
				}
			}
		}
	}()
	////
	// go func() {
	// 	for {
	// 		select {
	// 		case <-s.ctx.Done():
	// 			return
	// 		case txs := <-s.transferCh:
	// 			s.persistenceTransfer(txs)
	// 		}
	// 	}
	// }()
}
func (s *EthModule) updateHeight(number int64) {

	g.Log().Infof(s.ctx, "chainId:%d, updateHeight: %d", s.chainId, number)
	err := service.DB().UpdateState(s.ctx, s.chainId, number)
	if err != nil {
		g.Log().Fatalf(s.ctx, "fail to update height, err: %s", err)
	}
}

// //
// //
func (s *EthModule) Info() string {
	return fmt.Sprintf("%s|%d|%d|%d,contracts:%d", s.chainId, s.currentBlock, s.lastBlock, len(s.syncContracts))
}

func (s *EthModule) ChainId() int64 {
	return s.chainId
}
func (s *EthModule) LastBlock() int64 {
	return s.lastBlock
}
func (s *EthModule) syncBlock(latestBlock int64) error {

	topHeight := latestBlock - conf.Config.Syncing.WaitBlock
	g.Log().Infof(s.ctx, "chainId:%d, get header. latest: %d, topHeight: %d, current: %d, wait:%d", s.chainId, latestBlock, topHeight, s.currentBlock, conf.Config.Syncing.WaitBlock)
	////
	//// syncbatchblock
	startNumber := s.currentBlock + 1
	if topHeight > startNumber {
		endNumber := s.currentBlock + conf.Config.Syncing.BatchSyncTask
		if endNumber > topHeight {
			endNumber = topHeight
		}
		g.Log().Info(s.ctx, "syncBlock from:", startNumber, "end:", endNumber)
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
				txs, err := s.processBlock(s.ctx, blockNumber, s.cli)
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
			return errors.New("batchSync err")
		}
		///
		////send event
		txs := []*entity.SyncchainChainTransfer{}
		for _, tx := range txsmap {
			txs = append(txs, tx...)
		}
		if len(txs) > 0 {
			err := service.DB().UpTransactionMap(s.ctx, s.chainId, txsmap)
			if err != nil {
				g.Log().Fatal(s.ctx, "InsertTransfer_Transaction:", err)
			}
			for i, txs := range txsmap {
				service.EvnetSender().SendEvnetBatch_Latest(s.ctx, txs)
				g.Log().Debugf(s.ctx, "persistenceTransfer cached,chainId:%d , number:%d, log:%d", s.chainId, i, len(txs))
			}
			service.EvnetSender().SendEvnetBatch(s.ctx, txs)
		}

		g.Log().Infof(s.ctx, "%d:syncBlock, startNumber: %d, endNumber: %d, cnt:%d", s.chainId, startNumber, endNumber, len(txs))
		s.currentBlock = endNumber
	} else {
		return errors.New("syncBlock err")
	}
	return nil
}
