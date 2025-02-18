package conf

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

type Cache struct {
	SessionDuration int `json:"sessionDuration" v:"required|min:100"`
}
type Syncing struct {
	RpcUrl        string `json:"rpcUrl" v:"required"`
	MsgSize       int64  `json:"msgSize" v:"required|min:100000"`
	BatchSyncTask int64  `json:"batchSyncTask" v:"required|min:1"`
	BlockInterval int64  `json:"blockInterval" v:"required|min:1"`
	TimeOut       int    `json:"timeOut" v:"required|min:1"`
	WaitBlock     int64  `json:"waitBlock" v:"required|min:1"`
}
type Server struct {
	Address string `json:"address" v:"required"`
	WorkId  int    `json:"workId" v:"required|min:1"`
	Name    string `json:"name" v:"required"`
}
type Nrpcfg struct {
	NatsUrl string `json:"natsUrl" v:"required"`
}

// //

type Token2Native struct {
	ChainId  int64  `json:"chainId"`
	Contract string `json:"contract"`
}
type SkipAddrs struct {
	ChainId int64            `json:"chainId"`
	Address []common.Address `json:"contract"`
}
type Cfg struct {
	Server    *Server  `json:"server" v:"required"`
	Syncing   *Syncing `json:"syncing" v:"required"`
	Cache     *Cache   `json:"cache" v:"required"`
	JaegerUrl string   `json:"jaegerUrl" `
	Nrpc      *Nrpcfg  `json:"nrpc" v:"required"`
	////
	SyncCfgFile string `json:"syncCfgFile" v:"required"`
}

var Config = &Cfg{}

func init() {
	ctx := gctx.GetInitCtx()
	cfg := gcfg.Instance()
	v, err := cfg.Data(ctx)
	if err != nil {
		panic(err)
	}
	val := gvar.New(v)
	err = val.Structs(Config)
	if err != nil {
		panic(err)
	}
	if err := g.Validator().Data(Config).Run(ctx); err != nil {
		panic(err)
	}
	/////

}
