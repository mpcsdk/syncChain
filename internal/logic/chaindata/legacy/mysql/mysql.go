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
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	mysqlDBLog *sql.DB
	mysqlErr   error
)

func InitMySql() {
	mkWorkingDir()
	dsn := fmt.Sprintf("file:storage/chaindata.db?mode=rwc&_journal_mode=WAL&_cache_size=-500000")
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE if NOT EXISTS 'chaindata'(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		'chainid' INTEGER NOT NULL, 
		'height' INTEGER NOT NULL, 
		'blockhash' varchar(66) NOT NULL, 
		'ts' varchar(66) NOT NULL, 
		'txhash' varchar(66) NOT NULL, 
		'fromaddr' varchar(66) NOT NULL, 
		'toaddr' varchar(66) NOT NULL,
		'value' varchar(66) NOT NULL, 
		'contract' varchar(66) DEFAULT '',
		'gas' varchar(40) DEFAULT '',
		'gasprice' varchar(40) DEFAULT '', 
		UNIQUE ('chainid','txhash'));`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE INDEX if NOT EXISTS chainid ON chaindata (chainid);")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE INDEX if NOT EXISTS fromaddr ON chaindata (fromaddr);")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE INDEX if NOT EXISTS toddr ON chaindata (toaddr);")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE INDEX if NOT EXISTS contract ON chaindata (contract);")
	if err != nil {
		panic(err)
	}

	mysqlDBLog = db

	mysqlDBLog.SetMaxOpenConns(5)
	mysqlDBLog.SetMaxIdleConns(5)
	mysqlDBLog.SetConnMaxLifetime(100 * time.Second)

	if mysqlErr = mysqlDBLog.Ping(); nil != mysqlErr {
		mysqlDBLog.Close()
		panic(mysqlErr.Error())
	}
}

func mkWorkingDir() {
	path := "storage"
	_, err := os.Stat(path)
	if err == nil {
		return
	}

	os.MkdirAll(path, os.ModePerm)
}

func CloseMysql() {
	if nil != mysqlDBLog {
		mysqlDBLog.Close()
	}

}
