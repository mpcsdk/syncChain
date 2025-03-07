package tracetx

import (
	"context"
	"syncChain/internal/logic/chaindata/types"

	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type Empty struct {
}

func (s *Empty) GetTraceTransfer(ctx context.Context, block *types.Block) ([]*entity.SyncchainChainTransfer, error) {
	return nil, nil
}
