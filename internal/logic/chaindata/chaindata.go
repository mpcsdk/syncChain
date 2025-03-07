package chaindata

import (
	"context"
	block "syncChain/internal/logic/chaindata/sync"
	"syncChain/internal/service"

	"syncChain/internal/conf"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mq"
)

type sChainData struct {
	ctx         context.Context
	cancle      context.CancelFunc
	chainclient *block.EthModule
	closed      bool
	nats        *mq.NatsServer
	/////
	riskCtrlRule *mpcdao.RiskAdminDB
	// chainCfg     *mpcdao.ChainCfg
	////
}

//	func (s *sChainData) Close() {
//		s.closed = true
//		// s.chainclient.Close()
//	}
func (s *sChainData) ClientState() map[string]interface{} {
	d := map[string]interface{}{}
	d["chainId"] = s.chainclient.ChainId()
	d["latest"] = s.chainclient.LastBlock()
	d["confirmed"] = s.chainclient.LastBlock()

	return d
}

// //
type SkipAddrs struct {
	ChainId   int64    `json:"chainId"`
	Contracts []string `json:"contracts"`
}
type SyncCfg struct {
	SkipToAddr   []SkipAddrs `json:"skipToAddr"`
	SkipFromAddr []SkipAddrs `json:"skipFromAddr"`
}

func syncCfg(ctx context.Context, chainId int64) (*SyncCfg, error) {
	////
	cfg := gcfg.Instance(conf.Config.SyncCfgFile)
	v, err := cfg.Data(ctx)
	if err != nil {
		return nil, err
	}
	val := gvar.New(v)
	syncCfg := &SyncCfg{}
	err = val.Structs(syncCfg)
	if err != nil {
		panic(err)
	}
	if err := g.Validator().Data(syncCfg).Run(ctx); err != nil {
		return nil, err
	}
	////
	return syncCfg, nil
}

// //
func New() *sChainData {
	///syncchain
	ctx, cancle := context.WithCancel(gctx.GetInitCtx())

	p, err := gcmd.Parse(g.MapStrBool{
		"s,sync": true,
	})
	if err != nil {
		panic(err)
	}
	//
	opts := p.GetOptAll()
	syncChainId := opts["sync"]
	////
	///sycn chain
	if syncChainId == "" {
		panic("sync chaindata")
	}
	chainId := gconv.Int64(syncChainId)
	if chainId == 0 {
		panic("chainId")
	}

	briefs, err := service.DB().GetContractAbiBriefs(ctx, chainId)
	if err != nil {
		panic(err)
	}
	contracts := []common.Address{}
	for _, brief := range briefs {
		contracts = append(contracts, common.HexToAddress(brief.ContractAddress))
	}

	///init transfer db
	err = service.DB().InitChainTransferDB(ctx, chainId)
	if err != nil {
		panic(err)
	}
	///sync cfg
	syncBlockCfg, err := syncCfg(ctx, chainId)
	if err != nil {
		panic(err)
	}
	///
	//////
	s := &sChainData{
		ctx:         ctx,
		cancle:      cancle,
		chainclient: &block.EthModule{},
		closed:      false,
		//////
	}
	////
	skipToAddrs := []common.Address{}
	for _, skiptoaddr := range syncBlockCfg.SkipToAddr {
		if skiptoaddr.ChainId != chainId {
			continue
		}
		for _, addr := range skiptoaddr.Contracts {
			skipToAddrs = append(skipToAddrs, common.HexToAddress(addr))
		}
		break
	}
	////
	skipFromAddrs := []common.Address{}
	for _, skipaddr := range syncBlockCfg.SkipFromAddr {
		if skipaddr.ChainId != chainId {
			continue
		}
		for _, addr := range skipaddr.Contracts {
			skipFromAddrs = append(skipFromAddrs, common.HexToAddress(addr))
		}
		break
	}
	//////
	state, err := service.DB().GetState(ctx, chainId)
	if err != nil {
		panic(err)
	}
	//////
	module := block.NewEthModule(s.ctx,
		chainId,
		state.CurrentBlock,
		conf.Config.Syncing.RpcUrl, contracts, skipToAddrs, skipFromAddrs)
	//////
	module.Start()
	////
	s.chainclient = module

	return s
}
func (s *sChainData) Stop() {
	s.chainclient.Exit()
}
func (s *sChainData) IsRunning() bool {
	return s.chainclient.IsRunning()
}
