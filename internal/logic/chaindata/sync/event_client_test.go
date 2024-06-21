package block

import (
	"os"
	"syncChain/internal/logic/chaindata/common"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	defer func() {
		os.RemoveAll("logs")
	}()

	common.InitConf("chain.ini")

	// Init()

	time.Sleep(10 * time.Hour)
}
