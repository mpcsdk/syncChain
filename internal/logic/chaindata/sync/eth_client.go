package block

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

const (
	blockWait  = 10 * time.Second
	clientWait = 3 * time.Second
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
	exit   chan bool
	pause  chan bool
	closed bool
	///
	rpcList []string
	name    string
	chainId int64

	contracts contracts

	client *util.Client

	// last block from client
	count byte

	// last block processed
	lastBlock      int64
	confirmedBlock int64
	headerBlock    int64

	blockTimer  *time.Timer
	clientTimer *time.Timer

	// abi           abi.ABI
	// event         abi.Event
	// transferTopic string

	lock sync.Mutex

	//
	logger     *glog.Logger
	chaincfgdb *mpcdao.ChainCfg
	////
	transferCh     chan []*entity.ChainTransfer
	blockTransfers map[int64][]*entity.ChainTransfer
}

func NewEthModule(ctx context.Context, chainid int64, name string, rpcList []string, heigh int64, logger *glog.Logger) *EthModule {
	s := &EthModule{
		ctx:            ctx,
		chainId:        chainid,
		name:           name,
		lastBlock:      heigh,
		confirmedBlock: heigh,
		rpcList:        rpcList,
		logger:         logger,
		exit:           make(chan bool),
		pause:          make(chan bool),
		closed:         false,
		contracts: contracts{
			addresses: []common.Address{},
			names:     map[string]string{},
		},
		chaincfgdb: mpcdao.NewChainCfg(nil, 0),
		////
		transferCh:     make(chan []*entity.ChainTransfer, 100),
		blockTransfers: map[int64][]*entity.ChainTransfer{},
	}
	////
	s.blockTimer = time.NewTimer(2 * time.Second)
	s.clientTimer = time.NewTimer(1 * time.Second)
	s.clientTimer.Stop()
	s.blockTimer.Stop()
	// s.logger = log.GetLoggerByIndex(log.EVENT, s.name)
	// heightStr := common2.GlobalConf.GetString(chainsHeight, s.name, "0")

	///
	s.loop()
	//
	return s
}

func (s *EthModule) loop() {
	go func() {
		for {
			select {
			case <-s.clientTimer.C:
				func() {
					s.lock.Lock()
					defer s.lock.Unlock()

					s.logger.Warningf(s.ctx, "%s clientTimer getClient", s.name)
					s.getClient()
				}()
				break
			case <-s.blockTimer.C:
				s.syncBlock()
				break
			case p := <-s.pause:
				if p {
					s.logger.Notice(s.ctx, "pause:", s.name)
					s.blockTimer.Stop()
					s.clientTimer.Stop()
				} else {
					s.logger.Notice(s.ctx, "continue:", s.name)
					s.blockTimer.Reset(blockWait)
				}
			case <-s.exit:
				s.logger.Debugf(s.ctx, "exit, at height: %d", s.lastBlock)
				return
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

func (s *EthModule) getChainId() int64 {
	g.Log().Debug(s.ctx, "eth_getChainId:", s.chainId)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := s.client.ChainID(ctx)
	if err != nil {
		s.logger.Errorf(s.ctx, "%s fail to get chainId, err: %s, close client and reconnect", s.name, err)
		s.closeClient()
		return 0
	}

	// success, but no result
	if nil == id {
		s.logger.Errorf(s.ctx, "%s fail to get chainId, no id, close client and reconnect", s.name)
		s.closeClient()
	}

	chainId := id.Int64()
	s.logger.Warningf(s.ctx, "%s get chainId: %d", s.name, s.chainId)
	if 0 == chainId {
		s.logger.Errorf(s.ctx, "%s fail to get chainId, close client and reconnect", s.name)
	}
	return chainId
}

func (s *EthModule) getClient() *util.Client {
	if s.client != nil {
		return s.client
	}

	url := s.getURL()
	client, err := util.Dial(url)

	if err != nil {
		s.logger.Errorf(s.ctx, "fail to dial: %s", url)
		s.clientTimer.Reset(clientWait)
		return nil
	} else {
		s.logger.Infof(s.ctx, "dialed: %s", url)
	}

	s.client = client
	return client
}

func (s *EthModule) getURL() string {
	index := time.Now().Second() % len(s.rpcList)
	return strings.TrimSpace(s.rpcList[index])
}

func (s *EthModule) closeClient() {
	defer func() {
		if nil != s.clientTimer {
			s.clientTimer.Reset(clientWait)
		}

		s.count = 0
	}()

	if s.client == nil {
		return
	}

	s.client.Close()
	s.client = nil
}

func (s *EthModule) getHeader(client *util.Client) *types.Header {
	g.Log().Debug(s.ctx, "eth_Header:", s.chainId)
	var (
		header *types.Header
		err    error
	)
	ch := make(chan byte, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		header, err = client.HeaderByNumber(ctx, nil)
		ch <- 0
	}()

	select {
	case <-ch:
		if err != nil {
			s.logger.Error(s.ctx, "fail to get blockHeader:", s.chainId, "err:", err)
			s.closeClient()
			return nil
		}
		return header
	case <-ctx.Done():
		s.logger.Error(s.ctx, "fail to get blockHeader:", s.chainId, " timeout")
		s.closeClient()
		return nil
	}

}

func (s *EthModule) updateHeight(number int64) {
	if s.confirmedBlock >= number {
		return
	}
	s.confirmedBlock = number
	s.logger.Infof(s.ctx, "chainId:%d, updateHeight: %d", s.chainId, s.confirmedBlock)

	err := s.chaincfgdb.UpdateHeigh(s.ctx, s.chainId, s.confirmedBlock)
	if err != nil {
		s.logger.Errorf(s.ctx, "fail to update height, err: %s", err)
	}
}

// //
// //
func (s *EthModule) Info() string {
	return fmt.Sprintf("%s|%d|%d|%d,contracts:%d", s.name, s.chainId, s.confirmedBlock, s.lastBlock, s.contracts.Len())
}
func (s *EthModule) Close() {
	if s.closed {
		return
	}
	s.closed = true
	s.closeClient()
	s.exit <- true
}
func (s *EthModule) Pause() {
	if s.closed {
		return
	}

	s.pause <- true
}
func (s *EthModule) Continue() {
	if s.closed {
		return
	}

	s.pause <- false
}

// /
func (s *EthModule) ChainId() int64 {
	return s.chainId
}
func (s *EthModule) LastBlock() int64 {
	return s.lastBlock
}

// /
func (s *EthModule) Start() {
	s.blockTimer.Reset(blockWait)
}
func (s *EthModule) UpdateRpc(rpcs string) {
	s.rpcList = strings.Split(rpcs, ",")
}

// /
func (s *EthModule) UpdateContract(addr common.Address, name string) {
	s.contracts.Add(addr, name)
}
func (s *EthModule) DelContract(addr common.Address) {
	s.contracts.Del(addr)
	// for i, c := range s.contracts {
	// 	if c.Address == contract.Address {
	// 		s.contracts = append(s.contracts[:i], s.contracts[i+1:]...)
	// 		break
	// 	}
	// }
}
