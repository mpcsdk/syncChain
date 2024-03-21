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
	FromAddr string `json:"fromAddr"`
	ToAddr   string `json:"toAddr"`
	Contract string `json:"contract"`
	///
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	//
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
type QueryRes struct {
	g.Meta `mime:"text/html" example:"string"`
	Result interface{} `json:"result"`
}

// //
type StateReq struct {
	g.Meta `path:"/state" tags:"state" method:"post" summary:"You first hello api"`
}
type StateRes struct {
	g.Meta `mime:"text/html" example:"string"`
	Result interface{} `json:"result"`
}
