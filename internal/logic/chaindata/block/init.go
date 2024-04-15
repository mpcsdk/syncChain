package block

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var transferabi abi.ABI
var transferevent abi.Event
var transferTopic string

func init() {
	a, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		panic(err)
	}
	transferabi = a
	transferevent = transferabi.Events[transferName]
	transferTopic = hexutil.Encode(crypto.Keccak256([]byte(transferevent.Sig)))
	////
}
