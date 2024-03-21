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
	"os"
	"testing"
)

func TestInitMySql(t *testing.T) {
	defer func() {
		os.RemoveAll("logs")
		os.RemoveAll("logs-0.db")
		os.RemoveAll("logs-0.db-shm")
		os.RemoveAll("logs-0.db-wal")
		os.RemoveAll("1.ini")
		os.RemoveAll("storage0")
	}()

	InitMySql()

	InsertLogs(1, "2", "3", "4", "5", "6", "7", "8", "9", "10", "11")
}
