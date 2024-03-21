package query

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"syncChain/api/query/v1"
)

func (c *ControllerV1) Count(ctx context.Context, req *v1.CountReq) (res *v1.CountRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
