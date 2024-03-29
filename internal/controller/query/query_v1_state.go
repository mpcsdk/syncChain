package query

import (
	"context"
	v1 "syncChain/api/query/v1"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) State(ctx context.Context, req *v1.StateReq) (res *v1.StateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
