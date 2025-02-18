package syncBlock

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"syncChain/internal/conf"
	tracetx "syncChain/internal/logic/chaindata/sync/traceTx"
	"syncChain/internal/service"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type contract struct {
	Address common.Address
	Name    string
}
type contracts struct {
	addresses []common.Address
	names     map[string]string
}

func (s *contracts) Name(addr string) string {
	return s.names[addr]
}
func (s *contracts) Len() int {
	return len(s.addresses)
}
func (s *contracts) Add(addr common.Address, name string) {
	s.addresses = append(s.addresses, addr)
	s.names[addr.Hex()] = name
}
func (s *contracts) Del(addr common.Address) {
	for i, c := range s.addresses {
		if c.Cmp(addr) == 0 {
			s.addresses = append(s.addresses[:i], s.addresses[i+1:]...)
			delete(s.names, addr.Hex())
			break
		}
	}
}
func (s *contracts) Addresses() []common.Address {
	return s.addresses
}

type EthModule struct {
	ctx    context.Context
	cli    *ethclient.Client
	tracer tracetx.ITraceSyncer
	// exit   chan bool
	// pause  chan bool
	// closed bool
	///
	rpcList []string
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

	blockTimer *time.Timer
	blockWait  time.Duration
	// = time.Duration(conf.Config.Server.SyncInterval) * time.Second
	// clientTimer *time.Timer

	// abi           abi.ABI
	// event         abi.Event
	// transferTopic string
	lock sync.Mutex
	//
	// logger     *glog.Logger
	// chaincfgdb *mpcdao.ChainCfg
	////
	transferCh     chan []*entity.ChainTransfer
	blockTransfers map[int64][]*entity.ChainTransfer
	///
}

// var rpgtraceurl = "https://mainnet.rangersprotocol.com/api"
// var rpgtraceurl_testnet = "https://robin-api.rangersprotocol.com"

func NewEthModule(ctx context.Context, chainId int64, currentBlock int64, rpcList []string, syncContracts []common.Address, skipToAddrs []common.Address, skipFromAddrs []common.Address) *EthModule {
	cli, err := ethclient.Dial(rpcList[0])
	if err != nil {
		panic(err)
	}
	///
	tracer := tracetx.NewTraceSyncer(
		ctx, chainId,
		rpcList[0],
		time.Duration(conf.Config.Server.TimeOut)*time.Second,
	)
	s := &EthModule{
		ctx:          ctx,
		chainId:      chainId,
		currentBlock: currentBlock,
		// startBlock:   currentBlock,
		// lastBlock: currentBlock,
		rpcList: rpcList,
		cli:     cli,
		tracer:  tracer,
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
		transferCh:     make(chan []*entity.ChainTransfer, 100),
		blockTransfers: map[int64][]*entity.ChainTransfer{},
		///
		blockWait: time.Duration(conf.Config.Server.BlockInterval) * time.Second,
	}
	// if chainId == 9527 {
	// 	///rpgtestnet
	// 	cli, err := util.DialContext(ctx, rpgtraceurl_testnet)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	s.rpgtracecli = cli
	// } else if chainId == 2025 {
	// 	////rpg
	// 	cli, err := util.DialContext(ctx, rpgtraceurl)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	s.rpgtracecli = cli
	// }
	////
	s.blockTimer = time.NewTimer(s.blockWait)
	s.blockTimer.Stop()
	///
	s.loop()
	//
	return s
}

func (s *EthModule) loop() {
	go func() {
		s.blockTimer.Reset(s.blockWait)
		for {
			select {
			// case <-s.clientTimer.C:
			// 	func() {
			// 		s.lock.Lock()
			// 		defer s.lock.Unlock()

			// 		g.Log().Warningf(s.ctx, "%s clientTimer getClient", s.name)
			// 		s.getClient()
			// 	}()
			// 	break
			case <-s.ctx.Done():
				g.Log().Info(s.ctx, "exit, at height: ", s.currentBlock)
				return
			case <-s.blockTimer.C:
				s.syncBlock()
				break
				// case p := <-s.pause:
				// 	if p {
				// 		g.Log().Notice(s.ctx, "pause:", s.name)
				// 		s.blockTimer.Stop()
				// 		// s.clientTimer.Stop()
				// 	} else {
				// 		g.Log().Notice(s.ctx, "continue:", s.name)
				// 		s.blockTimer.Reset(s.blockWait)
				// 	}
				// case <-s.exit:
				// 	g.Log().Debugf(s.ctx, "exit, at height: %d", s.lastBlock)
				// 	return
				// }
			}
		}
	}()
	////
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case txs := <-s.transferCh:
				s.persistenceTransfer(txs)
			}
		}
	}()
}

// func (s *EthModule) getChainId() int64 {
// 	g.Log().Debug(s.ctx, "eth_getChainId:", s.chainId)
// 	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
// 	defer cancel()

// 	id, err := s.client.ChainID(ctx)
// 	if err != nil {
// 		g.Log().Errorf(s.ctx, "%s fail to get chainId, err: %s, close client and reconnect", s.name, err)
// 		s.closeClient()
// 		return 0
// 	}

// 	// success, but no result
// 	if nil == id {
// 		g.Log().Errorf(s.ctx, "%s fail to get chainId, no id, close client and reconnect", s.name)
// 		s.closeClient()
// 	}

// 	chainId := id.Int64()
// 	g.Log().Warningf(s.ctx, "%s get chainId: %d", s.name, s.chainId)
// 	if 0 == chainId {
// 		g.Log().Errorf(s.ctx, "%s fail to get chainId, close client and reconnect", s.name)
// 	}
// 	return chainId
// }

// func (s *EthModule) getClient() *util.Client {
// 	if s.client != nil {
// 		return s.client
// 	}

// 	url := s.getURL()
// 	client, err := util.Dial(url)

// 	if err != nil {
// 		g.Log().Errorf(s.ctx, "fail to dial: %s", url)
// 		// s.clientTimer.Reset(clientWait)
// 		return nil
// 	} else {
// 		g.Log().Infof(s.ctx, "dialed: %s", url)
// 	}

// 	s.client = client
// 	return client
// }

func (s *EthModule) getURL() string {
	index := time.Now().Second() % len(s.rpcList)
	return strings.TrimSpace(s.rpcList[index])
}

// func (s *EthModule) closeClient() {
// 	// defer func() {
// 	// 	if nil != s.clientTimer {
// 	// 		s.clientTimer.Reset(clientWait)
// 	// 	}

// 	// }()

// 	if s.client == nil {
// 		return
// 	}

// 	s.client.Close()
// 	s.client = nil
// }

func (s *EthModule) updateHeight(number int64) {

	g.Log().Infof(s.ctx, "chainId:%d, updateHeight: %d", s.chainId, number)
	err := service.DB().RiskAdmin().UpdateHeigh(s.ctx, s.chainId, number)
	if err != nil {
		g.Log().Fatalf(s.ctx, "fail to update height, err: %s", err)
	}
}

// //
// //
func (s *EthModule) Info() string {
	return fmt.Sprintf("%s|%d|%d|%d,contracts:%d", s.chainId, s.currentBlock, s.lastBlock, len(s.syncContracts))
}

//	func (s *EthModule) Close() {
//		if s.closed {
//			return
//		}
//		s.closed = true
//		s.closeClient()
//		s.exit <- true
//	}
// func (s *EthModule) Pause() {
// 	if s.closed {
// 		return
// 	}

//	}
// func (s *EthModule) Continue() {
// 	if s.closed {
// 		return
// 	}

// 	s.pause <- false
// }

// /
func (s *EthModule) ChainId() int64 {
	return s.chainId
}
func (s *EthModule) LastBlock() int64 {
	return s.lastBlock
}

// /
func (s *EthModule) Start() {
	s.blockTimer.Reset(s.blockWait)
}
func (s *EthModule) UpdateRpc(rpcs string) {
	s.rpcList = strings.Split(rpcs, ",")
}

// /
// func (s *EthModule) UpdateContract(addr common.Address, name string) {
// 	s.contracts.Add(addr, name)
// }
// func (s *EthModule) DelContract(addr common.Address) {
// 	s.contracts.Del(addr)
// 	// for i, c := range s.contracts {
// 	// 	if c.Address == contract.Address {
// 	// 		s.contracts = append(s.contracts[:i], s.contracts[i+1:]...)
// 	// 		break
// 	// 	}
// 	// }
// }
