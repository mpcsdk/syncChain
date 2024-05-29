package block

import (
	"context"
	"math/big"
	"sync"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
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

	s.logger.Debugf(s.ctx, "chainId:%d , start getting blocks. from %d to %d", s.chainId, s.lastBlock, last)
	for i := s.lastBlock + 1; i < last; i++ {
		block, _, txFroms, txHashes := s.getBlock(i, client)
		if nil == block {
			return
		}
		s.logger.Debugf(s.ctx, "getBlock,chainId:%d , block:%s,number:%d, txCount: %d", s.chainId, block.Hash().String(), i, len(block.Transactions()))

		// blockhashString := block.Hash().String()
		// if nil != blockhash {
		// 	blockhashString = blockhash.String()
		// }

		for index, tx := range block.Transactions() {
			receipt := s.getReceipt(tx.Hash(), client)
			////
			if nil == receipt {
				receipt = &types.Receipt{
					Status: types.ReceiptStatusFailed,
				}
			} else {
				if receipt.TxHash.Hex() != tx.Hash().Hex() {
					receipt.Status = types.ReceiptStatusFailed
				}
			}
			////
			s.processTx(block, tx, txFroms, txHashes, index, int64(receipt.Status))
			s.logger.Debugf(s.ctx, "getTransaction,chainId:%d , txHash:%s", s.chainId, tx.Hash())
			if 0 != s.contracts.Len() {
				logs := s.getLogs(i, client)
				if len(logs) > 0 {
					s.logger.Debugf(s.ctx, "getLogs,chainId:%d , log:%d", s.chainId, len(logs))
					s.processEvent(tx.Hash(), int64(block.Time()), logs, int64(receipt.Status))
				}
			}
		}
		/////event
		// for _, tx := range block.Transactions() {
		// 	if 0 != s.contracts.Len() {
		// 		s.processEvent(tx.Hash(), int64(block.Time()), client)
		// 	}
		// }
		////
		s.lastBlock = i
		s.updateHeight()
	}
}

func (s *EthModule) getBlock(i int64, client *util.Client) (*types.Block, *common.Hash, []*common.Address, []*common.Hash) {
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

func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}
