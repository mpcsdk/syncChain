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

var (
	rpgAddrByte = common.HexToAddress("0x71d9CFd1b7AdB1E8eb4c193CE6FFbe19B4aeE0dB").Bytes()
	rpgAddr     = common.HexToAddress("0x71d9CFd1b7AdB1E8eb4c193CE6FFbe19B4aeE0dB")
)

func Process20(ctx context.Context, chainId int64, ts int64, log *types.Log, status int64) error {
	fromAddr := common.BytesToAddress(log.Topics[1].Bytes())
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())
	//
	out, err := event20Transfer.Inputs.Unpack(log.Data)
	if err != nil {
		return err
	}
	value := out[0].(*big.Int)
	// if nil == value || 0 == value.Sign() {
	// 	s.logger.Info(s.ctx, "value is nil or zero")
	// 	return
	// }
	contractAddr := log.Address.String()
	if 0 == bytes.Compare(rpgAddrByte, log.Address.Bytes()) {
		contractAddr = ""
	}

	// service.DB().ChainLogs().Insert(gctx.GetInitCtx(), &entity.ChainLogs{
	// 	// ChainId:   chainId,
	// 	// Height:    i,
	// 	BlockHash: log.BlockHash.String(),
	// 	Ts:        ts,
	// 	TxHash:    log.TxHash.String(),
	// 	TxIdx:     int(log.TxIndex),
	// 	From:      fromAddr.String(),
	// 	To:        toAddr.String(),
	// 	Contract:  contractAddr,
	// 	Value:     value.String(),
	// 	Gas:       "0",
	// 	GasPrice:  "0",
	// 	LogIdx:    int(log.Index),
	// 	Nonce:     0,
	// 	Kind:      "erc20",

	// })
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
		Kind:      "erc20",
		Removed:   log.Removed,
		Status:    status,
	}
	return service.DB().InsertTransfer(ctx, data)
}
