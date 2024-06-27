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

var (
// rpgAddrByte = common.HexToAddress("0x71d9CFd1b7AdB1E8eb4c193CE6FFbe19B4aeE0dB").Bytes()
// rpgAddr     = common.HexToAddress("0x71d9CFd1b7AdB1E8eb4c193CE6FFbe19B4aeE0dB")
)

func Process20(ctx context.Context, chainId int64, ts int64, log *types.Log) *entity.ChainTransfer {
	fromAddr := common.BytesToAddress(log.Topics[1].Bytes())
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())
	///
	out, err := util.Event20Transfer.Inputs.Unpack(log.Data)
	if err != nil {
		g.Log().Error(ctx, "unpack err", err)
		return nil
	}
	value := out[0].(*big.Int)

	contractAddr := log.Address.String()
	kind := "erc20"

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
		Kind:      kind,
		Removed:   log.Removed,
		Status:    0,
		TraceTag:  "",
	}
	return data
}
