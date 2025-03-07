package transfer

import (
	"context"
	"math/big"
	"syncChain/internal/logic/chaindata/types"
	"syncChain/internal/logic/chaindata/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func Process1155Batch(ctx context.Context, chainId int64, ts int64, log *types.Log) []*entity.SyncchainChainTransfer {
	// operator := common.BytesToAddress(log.Topics[1].Bytes())
	fromAddr := common.BytesToAddress(log.Topics[1].Bytes())
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())
	//

	out, err := util.Event1155mul.Inputs.Unpack(log.Data)
	if err != nil {
		g.Log().Error(ctx, "Process1155Batch:", err)
		return nil
	}
	contractAddr := log.Address.String()
	kind := "1155"

	////
	tokenIds := out[0].([]*big.Int)
	vals := out[1].([]*big.Int)
	///
	datas := []*entity.SyncchainChainTransfer{}
	for j, v := range vals {
		t := tokenIds[j]
		datas = append(datas, &entity.SyncchainChainTransfer{
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
			Kind:      kind,
			TokenId:   t.String(),
			Removed:   log.Removed,
			Status:    0,
			TraceTag:  "",
		})
	}
	return datas
}

func Process1155Signal(ctx context.Context, chainId int64, ts int64, log *types.Log) *entity.SyncchainChainTransfer {
	// operator := common.BytesToAddress(log.Topics[1].Bytes())
	fromAddr := common.BytesToAddress(log.Topics[2].Bytes())
	toAddr := common.BytesToAddress(log.Topics[3].Bytes())
	//
	out, err := util.Event1155signal.Inputs.Unpack(log.Data)
	if err != nil {
		g.Log().Error(ctx, "Process1155Signal:", err)
		return nil
	}
	tokenId := out[0].(*big.Int)
	value := out[1].(*big.Int)
	// if nil == value || 0 == value.Sign() {
	// 	g.Log().Info(s.ctx, "value is nil or zero")
	// 	// return
	// }
	contractAddr := log.Address.String()
	kind := "erc1155"

	data := &entity.SyncchainChainTransfer{
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
		Kind:      kind,
		TokenId:   tokenId.String(),
		Removed:   log.Removed,
		Status:    0,
		TraceTag:  "",
	}
	return data
}
