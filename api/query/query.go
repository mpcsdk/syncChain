// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package query

import (
	"context"
	
	"syncChain/api/query/v1"
)

type IQueryV1 interface {
	Count(ctx context.Context, req *v1.CountReq) (res *v1.CountRes, err error)
	Query(ctx context.Context, req *v1.QueryReq) (res *v1.QueryRes, err error)
	State(ctx context.Context, req *v1.StateReq) (res *v1.StateRes, err error)
}


