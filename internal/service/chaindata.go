// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	IChainData interface {
		//	func (s *sChainData) Close() {
		//		s.closed = true
		//		// s.chainclient.Close()
		//	}
		ClientState() map[string]interface{}
		Stop()
		IsRunning() bool
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
