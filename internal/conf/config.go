package conf

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

type Cache struct {
	SessionDuration int `json:"sessionDuration" v:"required|min:100"`
}

type Server struct {
	Address       string `json:"address" v:"required"`
	WorkId        int    `json:"workId" v:"required|min:1"`
	Name          string `json:"name" v:"required"`
	MsgSize       int64  `json:'msgSize" v:"required|min:100000"`
	BatchSyncTask int64  `json:"batchSyncTask" v:"required|min:1"`
}
type Nrpcfg struct {
	NatsUrl string `json:"natsUrl" v:"required"`
}

// //

type Token2Native struct {
	ChainId  int64  `json:"chainId"`
	Contract string `json:"contract"`
}
type SkipToAddr struct {
	ChainId   int64    `json:"chainId"`
	Contracts []string `json:"contract"`
}
type Cfg struct {
	Server    *Server `json:"server" v:"required"`
	Cache     *Cache  `json:"cache" v:"required"`
	JaegerUrl string  `json:"jaegerUrl" `
	Nrpc      *Nrpcfg `json:"nrpc" v:"required"`
	// Token2Native      []*Token2Native               `json:"token2Native" v:"required"`
	// Token2NativeChain map[int64]string              `json:"token2NativeChain"`
	// SkipToAddr        []*SkipToAddr                 `json:"skipToAddr"`
	// SkipToAddrChain   map[int64]map[string]struct{} `json:"skipToAddrChain"`
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
	//////
	// Config.Token2NativeChain = make(map[int64]string)
	// for _, v := range Config.Token2Native {
	// 	Config.Token2NativeChain[v.ChainId] = v.Contract
	// }
	// Config.SkipToAddrChain = make(map[int64]map[string]struct{})
	// for _, v := range Config.SkipToAddr {

	// 	if _, ok := Config.SkipToAddrChain[v.ChainId]; !ok {
	// 		Config.SkipToAddrChain[v.ChainId] = make(map[string]struct{})
	// 	}
	// 	for _, vv := range v.Contracts {
	// 		Config.SkipToAddrChain[v.ChainId][vv] = struct{}{}
	// 	}
	// }

}
