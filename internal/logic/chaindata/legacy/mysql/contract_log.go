// Copyright 2020 The RangersProtocol Authors
// This file is part of the RocketProtocol library.
//
// The RangersProtocol library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The RangersProtocol library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the RangersProtocol library. If not, see <http://www.gnu.org/licenses/>.

package mysql

import (
	"strings"
)

func Count(addr, chainId string) uint64 {
	sql := "select count(*) FROM chaindata WHERE (fromaddr = ? and chainid = ?);"
	rows, err := mysqlDBLog.Query(sql, addr, chainId)
	if nil != err {
		// logger.Errorf("fail to count, %s, %s", addr, chainId)
		return 0
	}
	defer rows.Close()

	for rows.Next() {
		var result uint64
		err := rows.Scan(&result)
		if err != nil {
			// logger.Errorf("fail to scan, %s, %s", addr, chainId)
			return 0
		}
		return result
	}

	return 0
}

type QueryReq struct {
	From     string
	To       string
	ChainId  string
	Contract string
	///
	FromBlock int64
	ToBlock   int64
}

func (s *QueryReq) Validate() bool {
	if s.FromBlock == 0 && s.ToBlock == 0 {
		return false
	}
	if s.ChainId == "" && s.Contract == "" && s.From == "" && s.To == "" {
		return false
	}

	return true
}

func QueryAdv(req *QueryReq) []item {
	/////
	sql := "select chainid,height,blockhash,ts,txhash,fromaddr,toaddr,`value`,contract,gas,gasprice  FROM chaindata WHERE "
	args := []interface{}{}

	if req.From != "" {
		args = append(args, strings.ToLower(req.From))
		sql += " lower(fromaddr) = ? and "
	}
	if req.To != "" {
		args = append(args, strings.ToLower(req.To))
		sql += " lower(toaddr) = ? and "
	}
	if req.ChainId != "" {
		args = append(args, req.ChainId)
		sql += " chainid = ? and "
	}
	if req.Contract != "" {
		args = append(args, strings.ToLower(req.Contract))
		sql += " lower(contract) = ? and "
	}
	argsCnt := len(args)
	if argsCnt == 0 {
		return nil
	}
	//
	if req.FromBlock != 0 {
		args = append(args, req.FromBlock)
		sql += " height >= ? and "
	}
	if req.ToBlock != 0 {
		args = append(args, req.ToBlock)
		sql += " height <= ? "
	}
	///
	sql = strings.TrimSpace(sql)
	sql = strings.TrimRight(sql, "and")
	sql += ";"
	//
	rows, err := mysqlDBLog.Query(sql, args...)
	if nil != err {
		// logger.Error("fail to count", req, err)
		return nil
	}
	defer rows.Close()
	//
	result := make([]item, 0)
	for rows.Next() {
		var data item
		err := rows.Scan(&data.ChainId, &data.Height, &data.Blockhash, &data.Ts, &data.Txhash, &data.Fromaddr, &data.Toaddr, &data.Value, &data.Contract, &data.Gas, &data.Gasprice)
		if err != nil {
			// logger.Error("fail to count", req, err)
			return nil
		}

		result = append(result, data)
	}

	return result
}

func Query(from, to, chainId string, page, pageSize uint64) []item {
	args := []interface{}{}
	sql := "select height,blockhash,ts,txhash,fromaddr,toaddr,`value`,contract,gas,gasprice FROM chaindata WHERE "
	if from != "" {
		sql += " lower(fromaddr) = ? and "
		args = append(args, strings.ToLower(from))
	}
	if to != "" {
		sql += " lower(toaddr) = ? and "
		args = append(args, strings.ToLower(to))
	}
	if chainId != "" {
		sql += " chainid = ? and "
		args = append(args, chainId)
	}
	//
	sql = strings.TrimSpace(sql)
	sql = strings.TrimRight(sql, "and")
	///
	sql += " limit ?, ?;"
	args = append(args, page*pageSize)
	args = append(args, pageSize)
	///
	rows, err := mysqlDBLog.Query(sql, args...)
	if nil != err {
		// logger.Error("fail to count,", args)
		return nil
	}
	defer rows.Close()

	result := make([]item, 0)
	for rows.Next() {
		var data item
		err := rows.Scan(&data.Height, &data.Blockhash, &data.Ts, &data.Txhash, &data.Fromaddr, &data.Toaddr, &data.Value, &data.Contract, &data.Gas, &data.Gasprice)
		if err != nil {
			// logger.Errorf("fail to scan, ", args)
			return nil
		}

		result = append(result, data)
	}

	return result
}

func InsertLogs(height int64, chainid, blockhash, ts, txhash, fromaddr, toaddr, value, contract, gas, gasprice string) {
	stmt, err := mysqlDBLog.Prepare("replace INTO chaindata(chainid,height,blockhash,ts,txhash,fromaddr,toaddr,`value`, contract,gas,gasprice) values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		// logger.Errorf("fail to prepare. chainId: %s, height: %d, blockhash: %s, txhash: %s", chainid, height, blockhash, txhash)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(chainid, height, blockhash, ts, txhash, fromaddr, toaddr, value, contract, gas, gasprice)
	if err != nil {
		// logger.Errorf("fail to exec. chainId: %s, height: %d, blockhash: %s, txhash: %s", chainid, height, blockhash, txhash)
		return
	}

	// row, _ := result.RowsAffected()
	// lastId, _ := result.LastInsertId()
	// logger.Infof("inserted. chainId: %s, height: %d, blockhash: %s, txhash: %s, lines: %d, lastId: %d", chainid, height, blockhash, txhash, row, lastId)
}

type item struct {
	ChainId   string `json:"chainid"`
	Height    int    `json:"height"`
	Blockhash string `json:"blockhash"`
	Ts        string `json:"timestamp"`
	Txhash    string `json:"txhash"`
	Toaddr    string `json:"toaddr"`
	Fromaddr  string `json:"fromaddr"`
	Value     string `json:"value"`
	Contract  string `json:"contract"`
	Gas       string `json:"gas"`
	Gasprice  string `json:"gasprice"`
}
