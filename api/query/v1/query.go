package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CountReq struct {
	g.Meta `path:"/count" tags:"count" method:"post" summary:"You first hello api"`
}
type CountRes struct {
	g.Meta `mime:"text/html" example:"string"`
	Count  int `json:"count"`
}

// /
type QueryReq struct {
	g.Meta   `path:"/query" tags:"query" method:"post" summary:"You first hello api"`
	ChainId  int64  `json:"chainId"`
	From     string `json:"from"`
	To       string `json:"to"`
	Contract string `json:"contract"`
	Kind     string `json:"kind"`
	///
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	//
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
type QueryResult struct {
	ChainId   int64  `json:"chainId"`
	BlockHash string `json:"blockHash"`
	TxHash    string `json:"txHash"`
	Ts        int64  `json:"ts"`
	From      string `json:"from"`
	To        string `json:"to"`
	Contract  string `json:"contract"`
	Kind      string `json:"kind"`
	Value     string `json:"value"`
	TokenId   string `json:"tokenId"`
	Symbol  string `json:"symbol"`
	Status    int64  `json:"status"`
}
type QueryRes struct {
	g.Meta `mime:"text/html" example:"string"`
	Result []*QueryResult `json:"result"`
}

// //
type StateReq struct {
	g.Meta `path:"/state" tags:"state" method:"post" summary:"You first hello api"`
}
type StateRes struct {
	g.Meta `mime:"text/html" example:"string"`
	Result interface{} `json:"result"`
}

// //
type ContractReq struct {
	g.Meta   `path:"/contracts" tags:"state" method:"post" summary:"You first hello api"`
	ChainId  int64  `json:"chainId"`
	Contract string `json:"contract"`
}
type ContractResData struct {
	ChainId  int64  `json:"chainId"`
	Contract string `json:"contract"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Decimal  int    `json:"decimal"`
}

type ContractRes struct {
	g.Meta    `mime:"text/html" example:"string"`
	Contracts []*ContractResData `json:"contracts"`
}
