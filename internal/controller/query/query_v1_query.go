package query

import (
	"context"

	v1 "syncChain/api/query/v1"
	"syncChain/internal/model"
	"syncChain/internal/service"

	"github.com/ethereum/go-ethereum/common"
)

func (c *ControllerV1) Query(ctx context.Context, req *v1.QueryReq) (res *v1.QueryRes, err error) {
	res = &v1.QueryRes{}

	query := &model.QueryTx{
		From:     common.HexToAddress(req.From).String(),
		To:       common.HexToAddress(req.To).String(),
		Contract: common.HexToAddress(req.Contract).String(),
		///
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		///
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	result, err := service.DB().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	//
	res.Result = result
	return res, nil
}
