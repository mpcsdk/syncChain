package chaindata

import (
	"context"
	"strings"
	block "syncChain/internal/logic/chaindata/sync"
	"syncChain/internal/logic/chaindata/util"
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
	riskCtrlRule *mpcdao.RiskCtrlRule
	// chainCfg     *mpcdao.ChainCfg
	////
}

func (s *sChainData) Close() {
	s.closed = true
	s.chainclient.Close()
}
func (s *sChainData) ClientState() map[string]interface{} {
	d := map[string]interface{}{}
	d["chainId"] = s.chainclient.ChainId()
	d["latest"] = s.chainclient.LastBlock()
	d["confirmed"] = s.chainclient.LastBlock()

	return d
}

// //
type SkipToAddr struct {
	ChainId int64    `json:"chainId"`
	Address []string `json:"contracts"`
}
type SyncCfg struct {
	SkipToAddr []SkipToAddr `json:"skipToAddr`
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
	////chaincfg from db
	chainCfg, err := service.DB().ChainCfg().GetCfg(ctx, chainId)
	if err != nil {
		panic(err)
	}
	rpcs := strings.Split(chainCfg.Rpc, ",")
	if len(rpcs) == 0 {
		panic(chainCfg)
	}
	cli, err := util.Dial(rpcs[0])
	if err != nil {
		panic(err)
	}
	cliChainId, err := cli.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	if cliChainId.Int64() != chainCfg.ChainId {
		panic("clichanId!=cfgChainId:" + rpcs[0])
	}
	////
	///filter contracts
	riskCtrlRule := mpcdao.NewRiskCtrlRule(nil, 0)
	briefs, err := riskCtrlRule.GetContractAbiBriefs(ctx, chainCfg.ChainId, "")
	if err != nil {
		panic(err)
	}
	contracts := []common.Address{}
	for _, brief := range briefs {
		contracts = append(contracts, common.HexToAddress(brief.ContractAddress))
	}
	// nats := mq.New(conf.Config.Nrpc.NatsUrl)
	// _, err = nats.CreateOrUpdateStream(mq.JetStream_SyncChain, []string{mq.JetSub_SyncChain}, conf.Config.Server.MsgSize)
	// if err != nil {
	// 	panic(err)
	// }
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
		for _, addr := range skiptoaddr.Address {
			skipToAddrs = append(skipToAddrs, common.HexToAddress(addr))
		}
		break
	}
	module := block.NewEthModule(s.ctx,
		chainCfg.Coin,
		chainCfg.ChainId,
		chainCfg.Heigh,
		rpcs, contracts, skipToAddrs)
	//////
	module.Start()
	////
	s.chainclient = module

	return s
}
