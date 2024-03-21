// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ChainData is the golang structure of table chain_data for DAO operations like Where/Data.
type ChainData struct {
	g.Meta    `orm:"table:chain_data, do:true"`
	ChainId   interface{} //
	Height    interface{} //
	BlockHash interface{} //
	Ts        interface{} //
	TxHash    interface{} //
	TxIdx     interface{} //
	LogIdx    interface{} //
	FromAddr  interface{} //
	ToAddr    interface{} //
	Contract  interface{} //
	Value     interface{} //
	Gas       interface{} //
	GasPrice  interface{} //
	Nonce     interface{} //
}
