package util

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/text/gstr"
)

type TraceAction struct {
	CallType      string         `json:"callType"`
	From          common.Address `json:"from"`
	To            common.Address `json:"to"`
	Gas           string         `json:"gas"`
	Input         string         `json:"input"`
	Init          string         `json:"init"`
	Author        string         `json:"author"`
	Value         string         `json:"value"`
	RewardType    string         `json:"rewardType"`
	Address       string         `json:"address"`
	RefundAddress string         `json:"refundAddress"`
	Balance       string         `json:"balance"`
}
type TraceResult struct {
	GasUsed string `json:"gasUsed"`
	Output  string `json:"output`
	Address string `json:"address"`
	Code    string `json:"code"`
}
type Trace struct {
	Action              TraceAction `json:"action"`
	BlockHash           common.Hash `json:"blockHash"`
	BlockNumber         int64       `json:"blockNumber"`
	Result              TraceResult `json:"result"`
	Subtraces           int         `json:"subtraces"`
	TraceAddress        []int       `json:"traceAddress"`
	TransactionHash     common.Hash `json:"transactionHash"`
	TransactionPosition int         `json:"transactionPosition"`
	Type                string      `json:"type"`
	tag                 string
}

func (s *Trace) Tag() string {
	if s.tag == "" {
		s.tag = s.Action.CallType + "_" + gstr.JoinAny(s.TraceAddress, "_")
	}
	return s.tag
}

// //
type TraceRpg struct {
	BlockHeight  int64          `json:"blockHeight"`
	BlockHash    common.Hash    `json:"blockHash"`
	Depth        int            `json:"depth"`
	GasLimit     string         `json:"gasLimit"`
	ParentTxHash common.Hash    `json:"parentTxHash"`
	TxIndex      int            `json:"txIndex"`
	Source       common.Address `json:"source"`
	Target       common.Address `json:"target"`
	Time         time.Time      `json:"time"`
	Type         string         `json:"type"`
	Value        string         `json:"value"`
	TraceTag     string         `json:"traceTag"`
}

type RpgTraceResult struct {
	Data []*TraceRpg `json:"data"`
}

func (ec *Client) TraceBlock_rpg(ctx context.Context, number int64) ([]*TraceRpg, error) {
	var data *RpgTraceResult
	err := ec.c.CallContext(ctx, &data, "Rocket_getInternalTxByBlock", number)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data.Data, err
}

func (ec *Client) TraceBlock(ctx context.Context, number *big.Int) ([]*Trace, error) {
	var head []*Trace
	err := ec.c.CallContext(ctx, &head, "trace_block", toBlockNumArg(number))
	if err != nil {
		return nil, err
	}
	return head, err
}
