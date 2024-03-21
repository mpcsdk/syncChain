// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChainDataDao is the data access object for table chain_data.
type ChainDataDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns ChainDataColumns // columns contains all the column names of Table for convenient usage.
}

// ChainDataColumns defines and stores column names for table chain_data.
type ChainDataColumns struct {
	ChainId   string //
	Height    string //
	BlockHash string //
	Ts        string //
	TxHash    string //
	TxIdx     string //
	LogIdx    string //
	FromAddr  string //
	ToAddr    string //
	Contract  string //
	Value     string //
	Gas       string //
	GasPrice  string //
	Nonce     string //
}

// chainDataColumns holds the columns for table chain_data.
var chainDataColumns = ChainDataColumns{
	ChainId:   "chain_id",
	Height:    "height",
	BlockHash: "block_hash",
	Ts:        "ts",
	TxHash:    "tx_hash",
	TxIdx:     "tx_idx",
	LogIdx:    "log_idx",
	FromAddr:  "from_addr",
	ToAddr:    "to_addr",
	Contract:  "contract",
	Value:     "value",
	Gas:       "gas",
	GasPrice:  "gas_price",
	Nonce:     "nonce",
}

// NewChainDataDao creates and returns a new DAO object for table data access.
func NewChainDataDao() *ChainDataDao {
	return &ChainDataDao{
		group:   "default",
		table:   "chain_data",
		columns: chainDataColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChainDataDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChainDataDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChainDataDao) Columns() ChainDataColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChainDataDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChainDataDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChainDataDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
