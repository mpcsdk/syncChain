// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// ChainData is the golang structure for table chain_data.
type ChainData struct {
	ChainId   int64  `json:"chainId"   ` //
	Height    int64  `json:"height"    ` //
	BlockHash string `json:"blockHash" ` //
	Ts        int64  `json:"ts"        ` //
	TxHash    string `json:"txHash"    ` //
	TxIdx     int    `json:"txIdx"     ` //
	LogIdx    int    `json:"logIdx"    ` //
	FromAddr  string `json:"fromAddr"  ` //
	ToAddr    string `json:"toAddr"    ` //
	Contract  string `json:"contract"  ` //
	Value     string `json:"value"     ` //
	Gas       string `json:"gas"       ` //
	GasPrice  string `json:"gasPrice"  ` //
	Nonce     int64  `json:"nonce"     ` //
}
