package block

import (
	"context"
	"math/big"
	"strconv"
	"sync"
	common2 "syncChain/internal/logic/chaindata/common"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/model/entity"
	"syncChain/internal/service"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogf/gf/v2/os/gctx"
	"golang.org/x/crypto/sha3"
)

var (
	big8 = big.NewInt(8)

	hasherPool = sync.Pool{
		New: func() interface{} {
			return sha3.NewLegacyKeccak256()
		},
	}
)

func (self *ethModule) processBlock() {
	self.lock.Lock()
	defer func() {
		self.blockTimer.Reset(blockWait)
		self.lock.Unlock()
	}()

	client := self.getClient()
	if nil == client {
		self.logger.Errorf(self.ctx, "fail to get client")
		return
	}

	header := self.getHeader(client)
	if nil == header {
		return
	}

	topHeight := header.Number.Int64()
	self.logger.Debugf(self.ctx, "get header. height: %d, hash: %s", topHeight, header.Hash().String())

	last := topHeight - 12
	if last == self.lastBlockFromClient {
		self.count++
		if self.count == 6 {
			self.logger.Warningf(self.ctx, "get max retry count, close client and reconnect")
			self.closeClient()
			return
		}
	} else {
		self.count = 0
		self.lastBlockFromClient = last
	}

	if last <= self.lastBlock {
		self.logger.Infof(self.ctx, "no need to processBlock, remote: %d, local: %d", last, self.lastBlock)
		return
	}

	self.logger.Debugf(self.ctx, "chainId:%d , start getting blocks. from %d to %d", self.chainId, self.lastBlock, last)
	for i := self.lastBlock + 1; i < last; i++ {
		block, blockhash, txFroms, txHashes := self.getBlock(i, client)
		if nil == block {
			return
		}

		blockhashString := block.Hash().String()
		if nil != blockhash {
			blockhashString = blockhash.String()
		}

		for index, tx := range block.Transactions() {
			value := tx.Value()
			if tx == nil || tx.To() == nil || 0 == value.Sign() {
				continue
			}

			toAddr := tx.To().String()
			gas := strconv.FormatUint(tx.Gas(), 10)
			gasPrice := tx.GasPrice().String()
			if nil != txFroms[index] {
				fromAddr := txFroms[index].String()
				txhash := txHashes[index].String()
				service.DB().Insert(gctx.GetInitCtx(), &entity.ChainData{
					ChainId:   tx.ChainId().Int64(),
					Height:    i,
					BlockHash: blockhashString,
					Ts:        int64(block.Time()),
					TxHash:    txhash,
					TxIdx:     index,
					FromAddr:  fromAddr,
					ToAddr:    toAddr,
					Contract:  "",
					Value:     tx.Value().String(),
					Gas:       gas,
					GasPrice:  gasPrice,
					LogIdx:    -1,
					Nonce:     int64(tx.Nonce()),
				})
				continue
			}

			var hash common.Hash
			v, r, s := tx.RawSignatureValues()
			V := v
			if tx.Protected() {
				V = new(big.Int).Sub(v, new(big.Int).Mul(tx.ChainId(), big.NewInt(2)))
				V.Sub(V, big8)

				hash = rlpHash([]interface{}{
					tx.Nonce(),
					tx.GasPrice(),
					tx.Gas(),
					tx.To(),
					tx.Value(),
					tx.Data(),
					tx.ChainId(), uint(0), uint(0),
				})
			} else {
				hash = rlpHash([]interface{}{
					tx.Nonce(),
					tx.GasPrice(),
					tx.Gas(),
					tx.To(),
					tx.Value(),
					tx.Data(),
				})
			}
			fromAddr, err := common2.RecoverPlain(hash, r, s, V)
			txHash := tx.Hash().String()
			if nil != err {
				self.logger.Errorf(self.ctx, "fail to calc fromAddr, txhash: %s", txHash)
			} else {
				service.DB().Insert(gctx.GetInitCtx(), &entity.ChainData{
					ChainId:   tx.ChainId().Int64(),
					Height:    i,
					BlockHash: blockhashString,
					Ts:        int64(block.Time()),
					TxHash:    txHash,
					TxIdx:     index,
					FromAddr:  fromAddr.String(),
					ToAddr:    toAddr,
					Contract:  "",
					Value:     tx.Value().String(),
					Gas:       gas,
					GasPrice:  gasPrice,
					LogIdx:    -1,
					Nonce:     int64(tx.Nonce()),
				})
			}

		}

		if 0 != len(self.list) {
			self.processEvent(i, blockhashString, int64(block.Time()), client)
		}

		self.lastBlock = i
		self.updateHeight()
	}
}

func (self *ethModule) getBlock(i int64, client *Client) (*types.Block, *common.Hash, []*common.Address, []*common.Hash) {
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
			self.logger.Errorf(self.ctx, "fail to get block, err: %s, close client and reconnect", err)
			self.closeClient()
			return nil, nil, nil, nil
		}

		return block, hash, txFroms, txHashes
	case <-ctx.Done():
		self.logger.Errorf(self.ctx, "fail to get blockHeader, err: timeout, close client and reconnect")
		self.closeClient()
		return nil, nil, nil, nil
	}
}

func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}
