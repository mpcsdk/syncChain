package block

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var transferabi abi.ABI
var transferevent abi.Event
var transferTopic string = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

var abi1155 abi.ABI
var event1155 abi.Event
var signalTopic string = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
var mulTopic string = "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb"

func init() {
	a, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		panic(err)
	}
	transferabi = a
	transferevent = transferabi.Events[transferName]
	transferTopic = hexutil.Encode(crypto.Keccak256([]byte(transferevent.Sig)))
	////
	a, err = abi.JSON(strings.NewReader(abiData1155))
	if err != nil {
		panic(err)
	}
	abi1155 = a
	event1155 = abi1155.Events[mulTransfer]
	signalTopic = hexutil.Encode(crypto.Keccak256([]byte(event1155.Sig)))
	mulTopic = hexutil.Encode(crypto.Keccak256([]byte(event1155.Sig)))
	//
}
