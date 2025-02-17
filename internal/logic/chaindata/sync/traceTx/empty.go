package tracetx

import (
	"context"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

type Empty struct {
}

func (s *Empty) GetTraceTransfer(ctx context.Context, block *ethtypes.Block) ([]*entity.ChainTransfer, error) {
	return nil, nil
}
