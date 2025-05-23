package transfer

import (
	"context"
	"syncChain/internal/logic/chaindata/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func Process721(ctx context.Context, chainId int64, ts int64, log *types.Log) *entity.SyncchainChainTransfer {

	fromAddr := common.BytesToAddress(log.Topics[1].Bytes())
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())
	tokenId := log.Topics[3].Big().String()
	//
	// out, err := event721.Inputs.Unpack(log.Data)
	// if err != nil {
	// 	g.Log().Debugf(s.ctx, "fail to unpack data.  err: %s", err)
	// 	return
	// }

	// value := out[0].(*big.Int)
	// if nil == value || 0 == value.Sign() {
	// 	g.Log().Info(s.ctx, "value is nil or zero")
	// 	return

	contractAddr := log.Address.String()
	kind := "erc721"

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
		Value:     "",
		TokenId:   tokenId,
		Gas:       "0",
		GasPrice:  "0",
		LogIdx:    int(log.Index),
		Nonce:     0,
		Kind:      kind,
		Removed:   log.Removed,
		Status:    0,
		TraceTag:  "",
	}
	return data
}
