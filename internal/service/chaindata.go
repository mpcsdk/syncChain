// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	IChainData interface {
		Close()
		ClientState() map[int64]int64
	}
)

var (
	localChainData IChainData
)

func ChainData() IChainData {
	if localChainData == nil {
		panic("implement not found for interface IChainData, forgot register?")
	}
	return localChainData
}

func RegisterChainData(i IChainData) {
	localChainData = i
}
