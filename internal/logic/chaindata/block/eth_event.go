package block

import (
	"bytes"
	"context"
	"math/big"
	"syncChain/internal/logic/chaindata/types"

	"syncChain/internal/service"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

const (
	transferName = "Transfer"
	abiData      = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"},{"name":"start","type":"uint256"},{"name":"phase","type":"uint256"},{"name":"duration","type":"uint256"},{"name":"revocable","type":"bool"}],"name":"mintVesting","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"tos","type":"address[]"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferTo","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"cap","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"vesting","type":"address"}],"name":"revokeVesting","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"unpause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"mint","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"value","type":"uint256"}],"name":"burn","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"isBurner","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"isPauser","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"},{"name":"name","type":"string"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"paused","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"renouncePauser","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"}],"name":"addPauser","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"pause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"}],"name":"addMinter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"renounceMinter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"isMinter","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"renounceBurner","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"}],"name":"addBurner","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"tos","type":"address[]"},{"name":"values","type":"uint256[]"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"name","type":"string"},{"name":"symbol","type":"string"},{"name":"decimals","type":"uint8"},{"name":"cap","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"name","type":"string"}],"name":"TransferExtend","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"PauserAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"PauserRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"BurnerAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"BurnerRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"MinterAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"MinterRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}]`
	///
	abiData721   = `[{"inputs":[{"internalType":"address","name":"_terraformsDataAddress","type":"address"},{"internalType":"address","name":"_terraformsAugmentationsAddress","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"approved","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"bool","name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Daydreaming","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":false,"internalType":"address","name":"terraformer","type":"address"}],"name":"Terraformed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"seed","type":"uint256"}],"name":"TokensRevealed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"},{"stateMutability":"payable","type":"fallback"},{"inputs":[],"name":"MAX_SUPPLY","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"OWNER_ALLOTMENT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"PRICE","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"REVEAL_TIMESTAMP","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"SUPPLY","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TOKEN_SCALE","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newAddress","type":"address"}],"name":"addTokenURIAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"approve","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"address","name":"authorizedDreamer","type":"address"}],"name":"authorizeDreamer","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"uint256[16]","name":"dream","type":"uint256[16]"}],"name":"commitDreamToCanvas","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"dreamers","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"numTokens","type":"uint256"}],"name":"earlyMint","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"earlyMintActive","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"enterDream","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"numTokens","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"mintingPaused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"numTokens","type":"uint256"}],"name":"ownerClaim","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"redeemMintpass","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"seed","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"bool","name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"mintpassHolders","type":"address[]"}],"name":"setMintpassHolders","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"setSeed","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256[]","name":"tokens","type":"uint256[]"},{"internalType":"uint256","name":"index","type":"uint256"}],"name":"setTokenURIAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"structureData","outputs":[{"components":[{"internalType":"uint256","name":"levelNumber","type":"uint256"},{"internalType":"uint256","name":"tokensOnLevel","type":"uint256"},{"internalType":"int256","name":"structureSpaceX","type":"int256"},{"internalType":"int256","name":"structureSpaceY","type":"int256"},{"internalType":"int256","name":"structureSpaceZ","type":"int256"}],"internalType":"struct Terraforms.StructureLevel[20]","name":"structure","type":"tuple[20]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"terraformsAugmentationsAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"toggleEarly","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"togglePause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"tokenByIndex","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenCharacters","outputs":[{"internalType":"string[32][32]","name":"","type":"string[32][32]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"tokenCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenHTML","outputs":[{"internalType":"string","name":"result","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenHeightmapIndices","outputs":[{"internalType":"uint256[32][32]","name":"","type":"uint256[32][32]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"index","type":"uint256"}],"name":"tokenOfOwnerByIndex","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenSVG","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenSupplementalData","outputs":[{"components":[{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"uint256","name":"level","type":"uint256"},{"internalType":"uint256","name":"xCoordinate","type":"uint256"},{"internalType":"uint256","name":"yCoordinate","type":"uint256"},{"internalType":"int256","name":"elevation","type":"int256"},{"internalType":"int256","name":"structureSpaceX","type":"int256"},{"internalType":"int256","name":"structureSpaceY","type":"int256"},{"internalType":"int256","name":"structureSpaceZ","type":"int256"},{"internalType":"string","name":"zoneName","type":"string"},{"internalType":"string[10]","name":"zoneColors","type":"string[10]"},{"internalType":"string[9]","name":"characterSet","type":"string[9]"}],"internalType":"struct Terraforms.TokenData","name":"result","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenTerrainValues","outputs":[{"internalType":"int256[32][32]","name":"","type":"int256[32][32]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"tokenToAuthorizedDreamer","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"tokenToCanvasData","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"tokenToDreamBlock","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"tokenToDreamer","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"tokenToPlacement","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"tokenToStatus","outputs":[{"internalType":"enum TerraformsDreaming.Status","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"internalType":"string","name":"result","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"tokenURIAddresses","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]`
	abiData1155  = `[{"inputs":[{"internalType":"string","name":"_baseURI","type":"string"},{"internalType":"uint256","name":"_maxID","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"MintingNotEnabled","type":"error"},{"inputs":[],"name":"TokenDoesNotExist","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"bool","name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"newBaseURI","type":"string"}],"name":"BaseURIUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bool","name":"newEnabled","type":"bool"}],"name":"EnabledUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"newMaxID","type":"uint256"}],"name":"MaxIDUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"user","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnerUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256[]","name":"ids","type":"uint256[]"},{"indexed":false,"internalType":"uint256[]","name":"amounts","type":"uint256[]"}],"name":"TransferBatch","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"TransferSingle","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"value","type":"string"},{"indexed":true,"internalType":"uint256","name":"id","type":"uint256"}],"name":"URI","type":"event"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address[]","name":"owners","type":"address[]"},{"internalType":"uint256[]","name":"ids","type":"uint256[]"}],"name":"balanceOfBatch","outputs":[{"internalType":"uint256[]","name":"balances","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"baseURI","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"enabled","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"isApprovedForAll","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxID","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256[]","name":"ids","type":"uint256[]"},{"internalType":"uint256[]","name":"amounts","type":"uint256[]"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"safeBatchTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"bool","name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"_baseURI","type":"string"}],"name":"setBaseURI","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bool","name":"_enabled","type":"bool"}],"name":"setEnabled","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_maxID","type":"uint256"}],"name":"setMaxID","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"setOwner","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"uri","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`
	signTransfer = "TransferSingle"
	mulTransfer  = "TransferBatch"
)

var (
	rpgAddrByte = common.HexToAddress("0x71d9CFd1b7AdB1E8eb4c193CE6FFbe19B4aeE0dB").Bytes()
	rpgAddr     = common.HexToAddress("0x71d9CFd1b7AdB1E8eb4c193CE6FFbe19B4aeE0dB")
)

func (self *EthModule) process721(i int64, blockhash string, ts int64, client *Client, log types.Log) {

	fromAddr := common.BytesToAddress(log.Topics[1].Bytes())
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())
	tokenId := log.Topics[3].Hex()
	//
	// out, err := event721.Inputs.Unpack(log.Data)
	// if err != nil {
	// 	self.logger.Debugf(self.ctx, "fail to unpack data.  err: %s", err)
	// 	return
	// }

	self.logger.Debugf(self.ctx, "chainId:%d,tx: %s, get transfer data: %v, from: %s, to: %s", self.chainId, log.TxHash.String(), log.Data, fromAddr.String(), toAddr.String())
	// value := out[0].(*big.Int)
	// if nil == value || 0 == value.Sign() {
	// 	self.logger.Info(self.ctx, "value is nil or zero")
	// 	return
	// }
	contractAddr := log.Address.String()
	if 0 == bytes.Compare(rpgAddrByte, log.Address.Bytes()) {
		contractAddr = ""
	}

	service.DB().ChainData().Insert(gctx.GetInitCtx(), &entity.ChainData{
		ChainId:   self.chainId,
		Height:    i,
		BlockHash: log.BlockHash.String(),
		Ts:        ts,
		TxHash:    log.TxHash.String(),
		TxIdx:     int(log.TxIndex),
		From:      fromAddr.String(),
		To:        toAddr.String(),
		Contract:  contractAddr,
		Value:     "",
		TokenId:   tokenId,
		Gas:       "0",
		GasPrice:  "0",
		LogIdx:    int(log.Index),
		Nonce:     0,
		Kind:      "721",
	})
}
func (self *EthModule) process1155Batch(i int64, blockhash string, ts int64, client *Client, log types.Log) {
	// operator := common.BytesToAddress(log.Topics[1].Bytes())
	fromAddr := common.BytesToAddress(log.Topics[1].Bytes())
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())
	//
	out, err := event1155mul.Inputs.Unpack(log.Data)
	if err != nil {
		self.logger.Debugf(self.ctx, "fail to unpack data.  err: %s", err)
		return
	}

	self.logger.Debugf(self.ctx, "chainId:%d,tx: %s, get transfer data: %v, from: %s, to: %s", self.chainId, log.TxHash.String(), out, fromAddr.String(), toAddr.String())
	// tokenId := out[0].(*big.Int)
	// value := out[1].(*big.Int)
	// if nil == value || 0 == value.Sign() {
	// 	self.logger.Info(self.ctx, "value is nil or zero")
	// 	// return
	// }
	contractAddr := log.Address.String()

	////
	pos := int64(0)
	idlen := out[pos].(*big.Int)
	pos++
	///
	tokenIds := []*big.Int{}
	for i := int64(0); i < idlen.Int64(); i++ {
		id := out[i+pos].(*big.Int)
		pos++
		tokenIds = append(tokenIds, id)
	}
	///
	vlen := out[pos].(*big.Int)
	pos++
	vals := []*big.Int{}
	for i := int64(0); i < vlen.Int64(); i++ {
		v := out[i+pos].(*big.Int)
		pos++
		vals = append(vals, v)
	}
	//
	if len(vals) != len(tokenIds) {
		self.logger.Error(self.ctx, "value and tokenId length not equal")
		return
	}
	for j, v := range vals {
		t := tokenIds[j]
		service.DB().ChainData().Insert(gctx.GetInitCtx(), &entity.ChainData{
			ChainId:   self.chainId,
			Height:    i,
			BlockHash: log.BlockHash.String(),
			Ts:        ts,
			TxHash:    log.TxHash.String(),
			TxIdx:     int(log.TxIndex),
			From:      fromAddr.String(),
			To:        toAddr.String(),
			Contract:  contractAddr,
			Value:     v.String(),
			Gas:       "0",
			GasPrice:  "0",
			LogIdx:    int(log.Index),
			Nonce:     0,
			Kind:      "1155",
			TokenId:   t.String(),
		})
	}

}

func (self *EthModule) process1155Signal(i int64, blockhash string, ts int64, client *Client, log types.Log) {
	// operator := common.BytesToAddress(log.Topics[1].Bytes())
	fromAddr := common.BytesToAddress(log.Topics[2].Bytes())
	toAddr := common.BytesToAddress(log.Topics[3].Bytes())
	//
	out, err := event1155signal.Inputs.Unpack(log.Data)
	if err != nil {
		self.logger.Debugf(self.ctx, "fail to unpack data.  err: %s", err)
		return
	}

	self.logger.Debugf(self.ctx, "chainId:%d,tx: %s, get transfer data: %v, from: %s, to: %s", self.chainId, log.TxHash.String(), out, fromAddr.String(), toAddr.String())
	tokenId := out[0].(*big.Int)
	value := out[1].(*big.Int)
	// if nil == value || 0 == value.Sign() {
	// 	self.logger.Info(self.ctx, "value is nil or zero")
	// 	// return
	// }
	contractAddr := log.Address.String()
	if 0 == bytes.Compare(rpgAddrByte, log.Address.Bytes()) {
		contractAddr = ""
	}

	service.DB().ChainData().Insert(gctx.GetInitCtx(), &entity.ChainData{
		ChainId:   self.chainId,
		Height:    i,
		BlockHash: log.BlockHash.String(),
		Ts:        ts,
		TxHash:    log.TxHash.String(),
		TxIdx:     int(log.TxIndex),
		From:      fromAddr.String(),
		To:        toAddr.String(),
		Contract:  contractAddr,
		Value:     value.String(),
		Gas:       "0",
		GasPrice:  "0",
		LogIdx:    int(log.Index),
		Nonce:     0,
		Kind:      "1155",
		TokenId:   tokenId.String(),
	})
}

func (self *EthModule) processTransfer(i int64, blockhash string, ts int64, client *Client, log types.Log) {
	fromAddr := common.BytesToAddress(log.Topics[1].Bytes())
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())
	//
	out, err := transferevent.Inputs.Unpack(log.Data)
	if err != nil {
		self.logger.Debugf(self.ctx, "fail to unpack data.  err: %s", err)
		return
	}

	self.logger.Debugf(self.ctx, "chainId:%d,tx: %s, get transfer data: %v, from: %s, to: %s", self.chainId, log.TxHash.String(), out, fromAddr.String(), toAddr.String())
	value := out[0].(*big.Int)
	if nil == value || 0 == value.Sign() {
		self.logger.Info(self.ctx, "value is nil or zero")
		return
	}
	contractAddr := log.Address.String()
	if 0 == bytes.Compare(rpgAddrByte, log.Address.Bytes()) {
		contractAddr = ""
	}

	service.DB().ChainData().Insert(gctx.GetInitCtx(), &entity.ChainData{
		ChainId:   self.chainId,
		Height:    i,
		BlockHash: log.BlockHash.String(),
		Ts:        ts,
		TxHash:    log.TxHash.String(),
		TxIdx:     int(log.TxIndex),
		From:      fromAddr.String(),
		To:        toAddr.String(),
		Contract:  contractAddr,
		Value:     value.String(),
		Gas:       "0",
		GasPrice:  "0",
		LogIdx:    int(log.Index),
		Nonce:     0,
		Kind:      "20",
	})
}
func (self *EthModule) processEvent(i int64, blockhash string, ts int64, client *Client) {
	logs := self.getLogs(i, client)
	if nil == logs {
		return
	}

	for _, log := range logs {
		topic := log.Topics[0].String()
		switch topic {
		case transferTopic:
			if len(log.Topics) == 3 {
				self.processTransfer(i, blockhash, ts, client, log)
			} else if len(log.Topics) == 4 {
				self.process721(i, blockhash, ts, client, log)
			} else {
				self.logger.Info(self.ctx, "unknown event topic: %s")
			}
		case signalTopic:
			self.process1155Signal(i, blockhash, ts, client, log)
		case mulTopic:
			self.process1155Batch(i, blockhash, ts, client, log)
		default:
			self.logger.Info(self.ctx, "unknown event topic: %s")
		}
	}
	// if log.Topics[0] != "transfer" {
	// 	continue
	// }

}

func (self *EthModule) getLogs(i int64, client *Client) []types.Log {
	var (
		logs []types.Log
		err  error
	)
	ch := make(chan byte, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		var query ethereum.FilterQuery
		query.FromBlock = big.NewInt(i)
		query.ToBlock = big.NewInt(i)
		query.Addresses = self.contracts
		logs, err = client.FilterLogs(ctx, query)

		ch <- 0
	}()

	select {
	case <-ch:
		if err != nil {
			self.logger.Errorf(self.ctx, "fail to get logs, err: %s, close client and reconnect", err)
			self.closeClient()
			return nil
		}

		// success, but no result
		if nil == logs {
			logs = []types.Log{}
		}
		return logs
	case <-ctx.Done():
		self.logger.Errorf(self.ctx, "fail to get logs, err: timeout, close client and reconnect")
		self.closeClient()
		return nil
	}
}
