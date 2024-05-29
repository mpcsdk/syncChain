package query

import (
	"context"
	"math"
	"math/big"

	v1 "syncChain/api/query/v1"
	"syncChain/internal/service"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpccode"
	"github.com/mpcsdk/mpcCommon/mpcdao"
)

func (c *ControllerV1) Query(ctx context.Context, req *v1.QueryReq) (res *v1.QueryRes, err error) {
	g.Log().Debug(ctx, "Query req:", req)
	///
	if req.From == "" && req.To == "" && req.Contract == "" {
		return nil, mpccode.CodeParamInvalid("from, to, contract can't be all empty")
	}
	if req.StartTime >= req.EndTime {
		return nil, mpccode.CodeParamInvalid("startTime >= endTime")
	}
	if req.Page < 0 || req.PageSize < 0 {
		return nil, mpccode.CodeParamInvalid("page or pageSize invalid")
	}
	///
	query := &mpcdao.QueryData{
		ChainId: req.ChainId,
		From: func() string {
			if req.From == "" {
				return ""
			} else {
				return common.HexToAddress(req.From).String()
			}
		}(),
		To: func() string {
			if req.To == "" {
				return ""
			} else {
				return common.HexToAddress(req.To).String()
			}
		}(),
		Contract: func() string {
			if req.Contract == "" {
				return ""
			} else {
				return common.HexToAddress(req.Contract).String()
			}
		}(),
		///
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		///
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	/////
	if req.Kind == "token" {
		query.Kinds = []string{"erc20", "external"}
	} else if req.Kind == "nft" {
		query.Kinds = []string{"erc721", "erc1155"}
	} else {
		query.Kinds = []string{"external", "erc20", "erc721", "erc1155"}
	}
	////
	results, err := service.DB().QueryTransfer(ctx, query)
	if err != nil {
		g.Log().Error(ctx, "Query err:", err)
		return nil, mpccode.CodeInternalError(mpccode.TraceId(ctx))
	}
	/////
	res = &v1.QueryRes{}
	for _, r := range results {

		res.Result = append(res.Result, &v1.QueryResult{
			ChainId:   r.ChainId,
			BlockHash: r.BlockHash,
			TxHash:    r.TxHash,
			Ts:        r.Ts,
			From:      r.From,
			To:        r.To,
			Contract:  r.Contract,
			Kind:      r.Kind,
			Status:    r.Status,
			Symbol: func() string {
				////
				if r.Kind == "external" {
					chain := c.chains[r.ChainId]
					if chain != nil {
						return chain.Coin
					}
				} else {
					contract := c.contracts[r.Contract]
					if contract != nil {
						return contract.ContractName
					}
				}
				return ""
			}(),
			Value: func() string {
				if r.Kind == "external" {
					fbalance := big.NewFloat(0)
					fbalance.SetString(r.Value)
					fval := fbalance.Quo(fbalance, big.NewFloat(math.Pow10(18)))
					s := fval.Text('f', -1)
					return s
				} else if r.Kind == "erc20" {
					contract := c.contracts[r.Contract]
					if contract != nil {
						fbalance := big.NewFloat(0)
						fbalance.SetString(r.Value)
						fval := fbalance.Quo(fbalance, big.NewFloat(math.Pow10(contract.Decimal)))
						return fval.Text('f', -1)
					} else {
						return r.Value
					}
				}
				return r.Value
			}(),
			TokenId: r.TokenId,
		})

	}
	return res, nil
}
