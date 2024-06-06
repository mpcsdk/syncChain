package block

import (
	"math/big"
	"strconv"
	common2 "syncChain/internal/logic/chaindata/common"
	"syncChain/internal/logic/chaindata/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func (s *EthModule) processTx(block *types.Block, tx *types.Transaction, txFroms []*common.Address, txHashes []*common.Hash, index int) *entity.ChainTransfer {
	value := tx.Value()
	if tx == nil || tx.To() == nil || 0 == value.Sign() {
		return nil
	}
	toAddr := tx.To().String()
	gas := strconv.FormatUint(tx.Gas(), 10)
	gasPrice := tx.GasPrice().String()
	if nil != txFroms[index] {
		fromAddr := txFroms[index].String()
		txhash := txHashes[index].String()
		data := &entity.ChainTransfer{
			ChainId:   s.chainId,
			Height:    block.Number().Int64(),
			BlockHash: block.Hash().Hex(),
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
			Status:    0,
			Removed:   false,
		}
		return data
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
			Height:    block.Number().Int64(),
			BlockHash: block.Hash().Hex(),
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
			Kind:      "external",
			Status:    0,
			Removed:   false,
		}
		return data
	}
	return nil
}
