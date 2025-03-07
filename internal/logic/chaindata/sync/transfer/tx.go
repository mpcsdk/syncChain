package transfer

import (
	"context"
	"math/big"
	"strconv"
	"sync"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"golang.org/x/crypto/sha3"
)

func ProcessBlock(ctx context.Context, chainId int64, block *types.Block, tx *types.Transaction, index int) *entity.SyncchainChainTransfer {
	value := tx.Value()
	if tx == nil || tx.To() == nil || 0 == value.Sign() {
		return nil
	}
	// if len(tx.Data()) > 0 {
	// 	return nil
	// }
	toAddr := tx.To().String()

	gas := strconv.FormatUint(tx.Gas(), 10)
	gasPrice := tx.GasPrice().String()
	// from, err := tx.Sender()
	from := tx.Sender()
	if from.Big().Sign() == 0 {
		v, r, s := tx.RawSignatureValues()

		fromAddr, err := util.RecoverPlain(tx.Hash(), r, s, v)
		if nil != err {
			g.Log().Error(ctx, "fail to calc fromAddr, err:", err)
		}
		from = fromAddr
	}
	////
	data := &entity.SyncchainChainTransfer{
		ChainId:   chainId,
		Height:    block.Number().Int64(),
		BlockHash: block.Hash().Hex(),
		Ts:        int64(block.Time()),
		TxHash:    tx.Hash().Hex(),
		TxIdx:     index,
		From:      from.Hex(),
		To:        toAddr,
		Contract:  common.Address{}.String(),
		Value:     tx.Value().String(),
		Gas:       gas,
		GasPrice:  gasPrice,
		LogIdx:    -1,
		Nonce:     int64(tx.Nonce()),
		Kind:      "external",
		Status:    0,
		Removed:   false,
		TraceTag:  "",
	}
	return data
}

// func ProcessTx(ctx context.Context, chainId int64, block *types.Block, tx *types.Transaction, txFroms []*common.Address, txHashes []*common.Hash, index int) *entity.SyncchainChainTransfer {
// 	value := tx.Value()
// 	if tx == nil || tx.To() == nil || 0 == value.Sign() {
// 		return nil
// 	}
// 	// if len(tx.Data()) > 0 {
// 	// 	return nil
// 	// }
// 	toAddr := tx.To().String()

// 	gas := strconv.FormatUint(tx.Gas(), 10)
// 	gasPrice := tx.GasPrice().String()
// 	if nil != txFroms[index] {
// 		fromAddr := txFroms[index].String()
// 		txhash := txHashes[index].String()
// 		data := &entity.SyncchainChainTransfer{
// 			ChainId:   chainId,
// 			Height:    block.Number().Int64(),
// 			BlockHash: block.Hash().Hex(),
// 			Ts:        int64(block.Time()),
// 			TxHash:    txhash,
// 			TxIdx:     index,
// 			From:      fromAddr,
// 			To:        toAddr,
// 			Contract:  common.Address{}.String(),
// 			Value:     tx.Value().String(),
// 			Gas:       gas,
// 			GasPrice:  gasPrice,
// 			LogIdx:    -1,
// 			Nonce:     int64(tx.Nonce()),
// 			Kind:      "external",
// 			Status:    0,
// 			Removed:   false,
// 			TraceTag:  "",
// 		}
// 		return data
// 	}

// 	var hash common.Hash
// 	v, r, S := tx.RawSignatureValues()
// 	V := v
// 	if tx.Protected() {
// 		V = new(big.Int).Sub(v, new(big.Int).Mul(tx.ChainId(), big.NewInt(2)))
// 		V.Sub(V, big8)

// 		hash = rlpHash([]interface{}{
// 			tx.Nonce(),
// 			tx.GasPrice(),
// 			tx.Gas(),
// 			tx.To(),
// 			tx.Value(),
// 			tx.Data(),
// 			tx.ChainId(), uint(0), uint(0),
// 		})
// 	} else {
// 		hash = rlpHash([]interface{}{
// 			tx.Nonce(),
// 			tx.GasPrice(),
// 			tx.Gas(),
// 			tx.To(),
// 			tx.Value(),
// 			tx.Data(),
// 		})
// 	}
// 	fromAddr, err := util.RecoverPlain(hash, r, S, V)
// 	txHash := tx.Hash().String()
// 	if nil != err {
// 		g.Log().Errorf(ctx, "fail to calc fromAddr, txhash: %s", txHash)
// 	} else {
// 		data := &entity.SyncchainChainTransfer{
// 			ChainId:   tx.ChainId().Int64(),
// 			Height:    block.Number().Int64(),
// 			BlockHash: block.Hash().Hex(),
// 			Ts:        int64(block.Time()),
// 			TxHash:    txHash,
// 			TxIdx:     index,
// 			From:      fromAddr.String(),
// 			To:        toAddr,
// 			Contract:  common.Address{}.String(),
// 			Value:     tx.Value().String(),
// 			Gas:       gas,
// 			GasPrice:  gasPrice,
// 			LogIdx:    -1,
// 			Nonce:     int64(tx.Nonce()),
// 			Kind:      "external",
// 			Status:    0,
// 			Removed:   false,
// 		}
// 		return data
// 	}
// 	return nil
// }

var (
	big8 = big.NewInt(8)

	hasherPool = sync.Pool{
		New: func() interface{} {
			return sha3.NewLegacyKeccak256()
		},
	}
)

func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}
