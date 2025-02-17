package syncBlock

// // /trace_block
// func (s *EthModule) intx(block *ethtypes.Block) ([]*entity.ChainTransfer, error) {
// 	transfers := []*entity.ChainTransfer{}
// 	if s.chainId == 9527 || s.chainId == 2025 {
// 		/// for rpg method
// 		traces, err := s.getTraceBlock_rpg(block.Number().Int64(), s.rpgtracecli)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tracetxs := transfer.ProcessInTxnsRpg(s.ctx, s.chainId, block, traces)
// 		if len(tracetxs) > 0 {
// 			transfers = append(transfers, tracetxs...)
// 		}
// 	} else if s.chainId == 5003 {
// 	} else if s.chainId == 5000 {
// 		//support mantle natvie
// 		traces, err := s.getDebug_TraceBlock(block.Number().Int64(), client)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tracetxs := transfer.ProcessInTxns_mantle(s.ctx, s.chainId, block, traces)
// 		if tracetxs != nil {
// 			transfers = append(transfers, tracetxs...)
// 		}
// 	} else if s.chainId == 1 ||
// 		s.chainId == 1115511 ||
// 		s.chainId == 56 ||
// 		s.chainId == 97 {
// 		///other chains
// 		traces, err := s.getTraceBlock(block.Number().Int64(), client)
// 		if err != nil {
// 			g.Log().Warning(s.ctx, "getTraceBlock:", "chainId:", s.chainId, "number:", blockNumber, "err:", err)
// 			//return nil, err
// 		} else {
// 			tracetxs := transfer.ProcessInTxns(s.ctx, s.chainId, block, traces)
// 			if tracetxs != nil {
// 				transfers = append(transfers, tracetxs...)
// 			}
// 		}
// 	} else {
// 		///no trace method
// 	}
// 	return transfers, nil
// }
