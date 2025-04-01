// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package query

import (
	"syncChain/api/query"
	"syncChain/internal/service"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type ControllerV1 struct{
	contracts map[string]*entity.RiskadminContractabi
	// chains map[int64]*entity.RiskadminChaincfg

}

func NewV1() query.IQueryV1 {
	s := &ControllerV1{
		contracts: make(map[string]*entity.RiskadminContractabi),
		// chains: make(map[int64]*entity.RiskadminChaincfg),
	}
	////
	contracts:= service.DB().RiskAdminRepo().AllContract()
	for _, c := range contracts {
		s.contracts[c.ContractAddress] = c
	}
	////
	// chains , err := service.DB().RiskAdmin().AllChainsCfg(ctx)
	// for _, c := range chains {
	// 	s.chains[c.ChainId] = c
	// }
	////
	return s
}

