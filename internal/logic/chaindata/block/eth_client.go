package block

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"syncChain/internal/logic/chaindata/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/mpcsdk/mpcCommon/mpcdao"
)

const (
	blockWait  = 10 * time.Second
	clientWait = 3 * time.Second
)

type EthModule struct {
	ctx    context.Context
	exit   chan bool
	pause  chan bool
	closed bool
	///
	rpcList []string
	name    string
	chainId int64

	contracts []common.Address

	client *Client

	// last block from client
	lastBlockFromClient int64
	count               byte

	// last block processed
	lastBlock int64

	blockTimer  *time.Timer
	clientTimer *time.Timer

	abi           abi.ABI
	event         abi.Event
	transferTopic string

	lock sync.Mutex

	//
	logger     *glog.Logger
	chaincfgdb *mpcdao.ChainCfg
}

// self.exit = make(chan byte)
// self.closed = false
// self.lock = sync.Mutex{}
// // self.logger = log.GetLoggerByIndex(log.EVENT, self.name)

// // heightStr := common2.GlobalConf.GetString(chainsHeight, self.name, "0")
// self.lastBlock = heigh

// self.list = make([]common.Address, 0)
// if self.name == "rpg" {
// 	self.list = append(self.list, rpgAddr)
// }

// ///
// self.list = addresses
// //

// self.blockTimer = time.NewTimer(2 * time.Second)
// self.clientTimer = time.NewTimer(1 * time.Second)
// self.clientTimer.Stop()

// self.loop()

func NewEthModule(ctx context.Context, name string, rpcList []string, addresses []common.Address, heigh int64, logger *glog.Logger) *EthModule {
	s := &EthModule{
		ctx:        ctx,
		name:       name,
		lastBlock:  heigh,
		rpcList:    rpcList,
		logger:     logger,
		exit:       make(chan bool),
		pause:      make(chan bool),
		closed:     false,
		contracts:  addresses,
		chaincfgdb: mpcdao.NewChainCfg(),
	}
	////
	s.blockTimer = time.NewTimer(2 * time.Second)
	s.clientTimer = time.NewTimer(1 * time.Second)
	s.clientTimer.Stop()
	s.blockTimer.Stop()
	// self.logger = log.GetLoggerByIndex(log.EVENT, self.name)
	// heightStr := common2.GlobalConf.GetString(chainsHeight, self.name, "0")

	///

	s.loop()
	//
	return s
}

func (self *EthModule) loop() {
	go func() {
		for {
			select {
			case <-self.clientTimer.C:
				func() {
					self.lock.Lock()
					defer self.lock.Unlock()

					self.logger.Warningf(self.ctx, "%s clientTimer getClient", self.name)
					self.getClient()
				}()
				break

			case <-self.blockTimer.C:
				self.processBlock()
				break
			case p := <-self.pause:
				if p {
					self.logger.Notice(self.ctx, "pause:", self.name)
					self.blockTimer.Stop()
					self.clientTimer.Stop()
				} else {
					self.logger.Notice(self.ctx, "continue:", self.name)
					self.blockTimer.Reset(blockWait)
				}
			case <-self.exit:
				self.logger.Debugf(self.ctx, "exit, at height: %d", self.lastBlock)
				return
			}
		}
	}()
}

func (self *EthModule) getChainId() int64 {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := self.client.ChainID(ctx)
	if err != nil {
		self.logger.Errorf(self.ctx, "%s fail to get chainId, err: %s, close client and reconnect", self.name, err)
		self.closeClient()
		return 0
	}

	// success, but no result
	if nil == id {
		self.logger.Errorf(self.ctx, "%s fail to get chainId, no id, close client and reconnect", self.name)
		self.closeClient()
	}

	chainId := id.Int64()
	self.logger.Warningf(self.ctx, "%s get chainId: %d", self.name, self.chainId)
	if 0 == chainId {
		self.logger.Errorf(self.ctx, "%s fail to get chainId, close client and reconnect", self.name)
	}
	return chainId
}

// func (self *EthModule) initChainId() {
// 	for {
// 		client := self.getClient()
// 		if nil == client {
// 			time.Sleep(1 * time.Second)
// 			continue
// 		}

// 		func() {
// 			ch := make(chan byte, 1)
// 			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()

// 			var (
// 				id  *big.Int
// 				err error
// 			)
// 			go func() {
// 				id, err = client.ChainID(ctx)
// 				ch <- 0
// 			}()

// 			select {
// 			case <-ch:
// 				if err != nil {
// 					self.logger.Errorf(self.ctx, "%s fail to get chainId, err: %s, close client and reconnect", self.name, err)
// 					self.closeClient()
// 					return
// 				}

// 				// success, but no result
// 				if nil == id {
// 					self.logger.Errorf(self.ctx, "%s fail to get chainId, no id, close client and reconnect", self.name)
// 					self.closeClient()
// 				}

// 				self.chainId = id.Int64()
// 				self.logger.Warningf(self.ctx, "%s get chainId: %d", self.name, self.chainId)
// 				return
// 			case <-ctx.Done():
// 				self.logger.Errorf(self.ctx, "%s fail to get logs, err: timeout, close client and reconnect", self.name)
// 				self.closeClient()
// 				return
// 			}
// 		}()

// 		if 0 != self.chainId {
// 			return
// 		}
// 		self.logger.Errorf(self.ctx, "%s fail to get chainId, close client and reconnect", self.name)
// 		time.Sleep(1 * time.Second)
// 	}

// }

func (self *EthModule) getClient() *Client {
	if self.client != nil {
		return self.client
	}

	url := self.getURL()
	client, err := Dial(url)

	if err != nil {
		self.logger.Errorf(self.ctx, "fail to dial: %s", url)
		self.clientTimer.Reset(clientWait)
		return nil
	} else {
		self.logger.Infof(self.ctx, "dialed: %s", url)
	}

	self.client = client
	return client
}

func (self *EthModule) getURL() string {
	index := time.Now().Second() % len(self.rpcList)
	return strings.TrimSpace(self.rpcList[index])
}

func (self *EthModule) closeClient() {
	defer func() {
		if nil != self.clientTimer {
			self.clientTimer.Reset(clientWait)
		}

		self.count = 0
	}()

	if self.client == nil {
		return
	}

	self.client.Close()
	self.client = nil
}

func (self *EthModule) getHeader(client *Client) *types.Header {
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
			self.logger.Errorf(self.ctx, "fail to get blockHeader, err: %s, close client and reconnect", err)
			self.closeClient()
			return nil
		}
		return header
	case <-ctx.Done():
		self.logger.Errorf(self.ctx, "fail to get blockHeader, err: timeout, close client and reconnect")
		self.closeClient()
		return nil
	}

}

func (self *EthModule) updateHeight() {
	self.logger.Infof(self.ctx, "chainId:%d, updateHeight: %d", self.chainId, self.lastBlock)

	err := self.chaincfgdb.UpdateHeigh(self.ctx, self.chainId, self.lastBlock)
	if err != nil {
		self.logger.Errorf(self.ctx, "fail to update height, err: %s", err)
	}
}

// //
// //
func (self *EthModule) Info() string {
	return fmt.Sprintf("%s|%d|%d,contracts:%d", self.name, self.chainId, self.lastBlock, len(self.contracts))
}
func (self *EthModule) Close() {
	if self.closed {
		return
	}
	self.closed = true
	self.closeClient()
	self.exit <- true
}
func (self *EthModule) Pause() {
	if self.closed {
		return
	}

	self.pause <- true
}
func (self *EthModule) Continue() {
	if self.closed {
		return
	}

	self.pause <- false
}

// /
func (self *EthModule) ChainId() int64 {
	return self.chainId
}
func (self *EthModule) LastBlock() int64 {
	return self.lastBlock
}

// /
func (self *EthModule) Start() {
	self.blockTimer.Reset(blockWait)
}
func (self *EthModule) UpdateRpc(rpcs string) {
	self.rpcList = strings.Split(rpcs, ",")
}

// /
func (self *EthModule) UpdateContract(contract common.Address) {
	self.contracts = append(self.contracts, contract)
}
func (self *EthModule) DelContract(contract common.Address) {
	for i, _ := range self.contracts {
		self.contracts = append(self.contracts[:i], self.contracts[i+1:]...)
		break
	}
}
