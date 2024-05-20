package event

import (
	"bytes"
	"context"
	"math/big"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/service"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func Process1155Batch(ctx context.Context, chainId int64, ts int64, log *types.Log, status int64) error {
	// operator := common.BytesToAddress(log.Topics[1].Bytes())
	fromAddr := common.BytesToAddress(log.Topics[1].Bytes())
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())
	//

	out, err := event1155mul.Inputs.Unpack(log.Data)
	if err != nil {
		return err
	}
	contractAddr := log.Address.String()

	////
	pos := int64(0)
	idlen := out[pos].(*big.Int)
	pos++
	///
	tokenIds := []*big.Int{}
	for i := int64(0); i < idlen.Int64(); i++ {
		id := out[i+pos].(*big.Int)
		pos++
		tokenIds = append(tokenIds, id)
	}
	///
	vlen := out[pos].(*big.Int)
	pos++
	vals := []*big.Int{}
	for i := int64(0); i < vlen.Int64(); i++ {
		v := out[i+pos].(*big.Int)
		pos++
		vals = append(vals, v)
	}
	//
	// if len(vals) != len(tokenIds) {
	// 	s.logger.Error(s.ctx, "value and tokenId length not equal")
	// 	return
	// }
	datas := []*entity.ChainTransfer{}
	for j, v := range vals {
		t := tokenIds[j]
		datas = append(datas, &entity.ChainTransfer{
			ChainId:   chainId,
			Height:    int64(log.BlockNumber),
			BlockHash: log.BlockHash.String(),
			Ts:        ts,
			TxHash:    log.TxHash.String(),
			TxIdx:     int(log.TxIndex),
			From:      fromAddr.String(),
			To:        toAddr.String(),
			Contract:  contractAddr,
			Value:     v.String(),
			Gas:       "0",
			GasPrice:  "0",
			LogIdx:    int(log.Index),
			Nonce:     0,
			Kind:      "erc1155",
			TokenId:   t.String(),
			Removed:   log.Removed,
			Status:    status,
		})
	}
	return service.DB().InsertTransferBatch(ctx, datas)
}

func Process1155Signal(ctx context.Context, chainId int64, ts int64, log *types.Log, status int64) error {
	// operator := common.BytesToAddress(log.Topics[1].Bytes())
	fromAddr := common.BytesToAddress(log.Topics[2].Bytes())
	toAddr := common.BytesToAddress(log.Topics[3].Bytes())
	//
	out, err := event1155signal.Inputs.Unpack(log.Data)
	if err != nil {
		return err
	}
	tokenId := out[0].(*big.Int)
	value := out[1].(*big.Int)
	// if nil == value || 0 == value.Sign() {
	// 	s.logger.Info(s.ctx, "value is nil or zero")
	// 	// return
	// }
	contractAddr := log.Address.String()
	if 0 == bytes.Compare(rpgAddrByte, log.Address.Bytes()) {
		contractAddr = ""
	}

	data := &entity.ChainTransfer{
		ChainId:   chainId,
		Height:    int64(log.BlockNumber),
		BlockHash: log.BlockHash.String(),
		Ts:        ts,
		TxHash:    log.TxHash.String(),
		TxIdx:     int(log.TxIndex),
		From:      fromAddr.String(),
		To:        toAddr.String(),
		Contract:  contractAddr,
		Value:     value.String(),
		Gas:       "0",
		GasPrice:  "0",
		LogIdx:    int(log.Index),
		Nonce:     0,
		Kind:      "erc1155",
		TokenId:   tokenId.String(),
		Removed:   log.Removed,
		Status:    status,
	}
	return service.DB().InsertTransfer(ctx, data)
}
