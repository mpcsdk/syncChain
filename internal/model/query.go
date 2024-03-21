package model

type QueryTx struct {
	FromAddr string `json:"fromAddr"`
	ToAddr   string `json:"toAddr"`
	Contract string `json:"contract"`
	///
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	///
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
