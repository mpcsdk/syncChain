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
	briefs, err := s.riskCtrlRule.GetContractRuleBriefs(s.ctx, data.ChainId, "")
	if err != nil {
		g.Log().Error(s.ctx, "GetContractRuleBriefs:", err)
	}
	constrats := []common.Address{}
	for _, v := range briefs {
		constrats = append(constrats, common.HexToAddress(v.ContractAddress))
	}
	///
	module := block.NewEthModule(s.ctx, data.Coin, rpcs, constrats, data.Heigh, g.Log("blocklog"))
	s.clients[data.Id] = module
	///
	if data.IsEnable == 1 {
		module.Start()
	}
}

// /
func (s *sChainData) updateOpt(data *entity.Chaincfg) {
	if v, ok := s.clients[data.Id]; ok {
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
	if v, ok := s.clients[data.Id]; ok {
		v.Close()
		delete(s.clients, data.Id)
	}
}

// /
func (s *sChainData) addOptContractRule(data *entity.Contractrule) {
	for _, v := range s.clients {
		if v.ChainId() == data.ChainId {
			v.UpdateContract(common.HexToAddress(data.ContractAddress))
			return

		}
	}
}

func (s *sChainData) updateOptContractRule(data *entity.Contractrule) {
	for _, v := range s.clients {
		if v.ChainId() == data.ChainId {
			v.UpdateContract(common.HexToAddress(data.ContractAddress))
			return
		}
	}

}
func (s *sChainData) deleteOptContractRule(data *entity.Contractrule) {
	for _, v := range s.clients {
		if v.ChainId() == data.ChainId {
			v.DelContract(common.HexToAddress(data.ContractAddress))
			return
		}
	}
}