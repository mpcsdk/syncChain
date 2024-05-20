package block

import (
	"context"
	"math/big"
	"strconv"
	"sync"
	common2 "syncChain/internal/logic/chaindata/common"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"
	"syncChain/internal/service"
	"time"

	"github.com/lib/pq"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
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

func (s *EthModule) processTx() {

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
	if s.lastBlock == 0 {
		s.lastBlock = topHeight
	}
	s.logger.Debugf(s.ctx, "get header. height: %d, hash: %s", topHeight, header.Hash().String())

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
		block, blockhash, txFroms, txHashes := s.getBlock(i, client)
		if nil == block {
			return
		}
		s.logger.Debugf(s.ctx, "getBlock,chainId:%d , block:%d, txCount: %d", s.chainId, i, len(block.Transactions()))

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
				data := &entity.ChainTransfer{
					ChainId:   tx.ChainId().Int64(),
					Height:    i,
					BlockHash: blockhashString,
					Ts:        int64(block.Time()),
					TxHash:    txhash,
					TxIdx:     index,
					From:      fromAddr,
					To:        toAddr,
					Contract:  "",
					Value:     tx.Value().String(),
					Gas:       gas,
					GasPrice:  gasPrice,
					LogIdx:    -1,
					Nonce:     int64(tx.Nonce()),
					Kind:      "external",
				}
				err := service.DB().InsertTransfer(gctx.GetInitCtx(), data)
				if err != nil {
					switch v := gerror.Cause(err).(type) {
					case *pq.Error:
						if v.Code == "23505" { // unique_violation
							g.Log().Warning(s.ctx, "duplicate tx, txhash: ", data)
						} else {
							g.Log().Fatal(s.ctx, "fail to insert tx, err: ", err)
						}
					default:
						g.Log().Fatal(s.ctx, "fail to insert tx, err: ", err)
					}
				}
				continue
			}

			var hash common.Hash
			v, r, S := tx.RawSignatureValues()
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
			fromAddr, err := common2.RecoverPlain(hash, r, S, V)
			txHash := tx.Hash().String()
			if nil != err {
				s.logger.Errorf(s.ctx, "fail to calc fromAddr, txhash: %s", txHash)
			} else {
				data := &entity.ChainTransfer{
					ChainId:   tx.ChainId().Int64(),
					Height:    i,
					BlockHash: blockhashString,
					Ts:        int64(block.Time()),
					TxHash:    txHash,
					TxIdx:     index,
					From:      fromAddr.String(),
					To:        toAddr,
					Contract:  "",
					Value:     tx.Value().String(),
					Gas:       gas,
					GasPrice:  gasPrice,
					LogIdx:    -1,
					Nonce:     int64(tx.Nonce()),
					Kind:      "",
				}
				err := service.DB().InsertTransfer(gctx.GetInitCtx(), data)
				if err != nil {
					switch v := err.(type) {
					case *pq.Error:
						if v.Code == "23505" { // unique_violation
							g.Log().Warning(s.ctx, "duplicate tx, txhash:", data)
						} else {
							g.Log().Fatal(s.ctx, "fail to insert tx, err: ", err)
						}
					default:
						g.Log().Fatal(s.ctx, "fail to insert tx, err: ", err)
					}
				}
			}

		}
		/////event
		for _, tx := range block.Transactions() {
			if 0 != s.contracts.Len() {
				s.processEvent(tx.Hash(), int64(block.Time()), client)
			}
		}
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
			s.logger.Errorf(s.ctx, "fail to get block, err: %s, close client and reconnect", err)
			s.closeClient()
			return nil, nil, nil, nil
		}

		return block, hash, txFroms, txHashes
	case <-ctx.Done():
		s.logger.Errorf(s.ctx, "fail to get blockHeader, err: timeout, close client and reconnect")
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
