package block

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	common2 "syncChain/internal/logic/chaindata/common"
	"syncChain/internal/logic/chaindata/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogf/gf/v2/os/glog"
)

const (
	blockWait  = 10 * time.Second
	clientWait = 3 * time.Second
)

type ethModule struct {
	ctx context.Context
	///
	rpcList []string
	name    string
	chainId int64

	list []common.Address

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

	exit   chan byte
	closed bool
	//
	logger *glog.Logger
}

func (self *ethModule) info() string {
	return fmt.Sprintf("%s|%d|%d", self.name, self.chainId, self.lastBlock)
}

func (self *ethModule) start(name, addresses string) {
	self.name = name

	abi, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		self.logger.Errorf(self.ctx, "fail to start. name: %s, err: %s", name, err)
		return
	}
	self.abi = abi
	self.event = abi.Events[transferName]
	self.transferTopic = hexutil.Encode(crypto.Keccak256([]byte(self.event.Sig)))

	self.exit = make(chan byte)
	self.closed = false
	self.lock = sync.Mutex{}
	// self.logger = log.GetLoggerByIndex(log.EVENT, self.name)

	heightStr := common2.GlobalConf.GetString(chainsHeight, self.name, "0")
	self.lastBlock, _ = strconv.ParseInt(heightStr, 10, 32)

	self.list = make([]common.Address, 0)
	if self.name == "rpg" {
		self.list = append(self.list, rpgAddr)
	}

	if 0 != len(addresses) {
		for _, address := range strings.Split(addresses, ",") {
			self.list = append(self.list, common.HexToAddress(strings.TrimSpace(address)))
		}
	}

	self.initChainId()

	self.blockTimer = time.NewTimer(2 * time.Second)
	self.clientTimer = time.NewTimer(1 * time.Second)
	self.clientTimer.Stop()

	self.loop()
}

func (self *ethModule) loop() {
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

			case <-self.exit:
				self.logger.Debugf(self.ctx, "exit, at height: %d", self.lastBlock)
				return
			}
		}
	}()
}

func (self *ethModule) close() {
	if self.closed {
		return
	}

	self.closeClient()
	self.closed = true
	self.exit <- 1

}

func (self *ethModule) initChainId() {
	for {
		client := self.getClient()
		if nil == client {
			time.Sleep(1 * time.Second)
			continue
		}

		func() {
			ch := make(chan byte, 1)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var (
				id  *big.Int
				err error
			)
			go func() {
				id, err = client.ChainID(ctx)
				ch <- 0
			}()

			select {
			case <-ch:
				if err != nil {
					self.logger.Errorf(self.ctx, "%s fail to get chainId, err: %s, close client and reconnect", self.name, err)
					self.closeClient()
					return
				}

				// success, but no result
				if nil == id {
					self.logger.Errorf(self.ctx, "%s fail to get chainId, no id, close client and reconnect", self.name)
					self.closeClient()
				}

				self.chainId = id.Int64()
				self.logger.Warningf(self.ctx, "%s get chainId: %d", self.name, self.chainId)
				return
			case <-ctx.Done():
				self.logger.Errorf(self.ctx, "%s fail to get logs, err: timeout, close client and reconnect", self.name)
				self.closeClient()
				return
			}
		}()

		if 0 != self.chainId {
			return
		}
		self.logger.Errorf(self.ctx, "%s fail to get chainId, close client and reconnect", self.name)
		time.Sleep(1 * time.Second)
	}

}

func (self *ethModule) getClient() *Client {
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

func (self *ethModule) getURL() string {
	index := time.Now().Second() % len(self.rpcList)
	return strings.TrimSpace(self.rpcList[index])
}

func (self *ethModule) closeClient() {
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

func (self *ethModule) getHeader(client *Client) *types.Header {
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

func (self *ethModule) updateHeight() {
	self.logger.Infof(self.ctx, "chainId:%d, updateHeight: %d", self.chainId, self.lastBlock)
	common2.GlobalConf.SetString(chainsHeight, self.name, strconv.FormatInt(self.lastBlock, 10))
}
