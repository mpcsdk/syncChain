package chaindata

import (
	"strings"
	"syncChain/internal/logic/chaindata/block"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func (s *sChainData) addOpt(data *entity.Chaincfg) {
	rpcs := strings.Split(data.Rpc, ",")
	///
	briefs, err := s.riskCtrlRule.GetContractAbiBriefs(s.ctx, data.ChainId, "")
	if err != nil {
		g.Log().Error(s.ctx, "GetContractRuleBriefs:", err)
	}
	///
	module := block.NewEthModule(s.ctx, data.ChainId, data.Coin, rpcs, data.Heigh, g.Log("blocklog"))
	for _, v := range briefs {
		module.UpdateContract(common.HexToAddress(v.ContractAddress), v.ContractName)
	}
	s.chainclients[data.ChainId] = module
	///
	if data.IsEnable == 1 {
		module.Start()
	}
}

// /
func (s *sChainData) updateOpt(data *entity.Chaincfg) {
	if v, ok := s.chainclients[data.ChainId]; ok {
		///
		if data.IsEnable == 0 {
			v.Pause()
		} else {
			v.Continue()
		}
		///
		if data.Rpc != "" {
			v.UpdateRpc(data.Rpc)
		}
	}

}
func (s *sChainData) deleteOpt(data *entity.Chaincfg) {
	if v, ok := s.chainclients[data.ChainId]; ok {
		v.Close()
		delete(s.chainclients, data.ChainId)
	}
}

// /
func (s *sChainData) addOptContractRule(data *entity.Contractrule) {
	for _, v := range s.chainclients {
		if v.ChainId() == data.ChainId {
			v.UpdateContract(common.HexToAddress(data.ContractAddress), data.ContractName)
			return
		}
	}
}

func (s *sChainData) updateOptContractRule(data *entity.Contractrule) {
	for _, v := range s.chainclients {
		if v.ChainId() == data.ChainId {
			v.UpdateContract(common.HexToAddress(data.ContractAddress), data.ContractName)
			return
		}
	}

}
func (s *sChainData) deleteOptContractRule(data *entity.Contractrule) {
	for _, v := range s.chainclients {
		if v.ChainId() == data.ChainId {
			v.DelContract(common.HexToAddress(data.ContractAddress))
			return
		}
	}
}
