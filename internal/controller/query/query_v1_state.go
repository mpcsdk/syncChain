package query

import (
	"context"
	v1 "syncChain/api/query/v1"
	"syncChain/internal/service"
)

func (c *ControllerV1) State(ctx context.Context, req *v1.StateReq) (res *v1.StateRes, err error) {
	stat := service.ChainData().ClientState()
	res = &v1.StateRes{
		Result: stat,
	}
	return res, nil
}
