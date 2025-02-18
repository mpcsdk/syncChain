package tracetx

// type EthTrace struct {
// 	TraceSyncer
// }

// func newEthTracer(ctx context.Context, chainId int64, url string, ctxTimeOut time.Duration) *EthTrace {
// 	cli, err := ethclient.Dial(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &EthTrace{
// 		TraceSyncer: TraceSyncer{
// 			cli:        cli,
// 			ctx:        ctx,
// 			ctxTimeOut: ctxTimeOut,
// 			chainId:    chainId,
// 		},
// 	}
// }

// type TraceAction struct {
// 	CallType      string         `json:"callType"`
// 	From          common.Address `json:"from"`
// 	To            common.Address `json:"to"`
// 	Gas           string         `json:"gas"`
// 	Input         string         `json:"input"`
// 	Init          string         `json:"init"`
// 	Author        string         `json:"author"`
// 	Value         *hexutil.Big   `json:"value"`
// 	RewardType    string         `json:"rewardType"`
// 	Address       string         `json:"address"`
// 	RefundAddress string         `json:"refundAddress"`
// 	Balance       string         `json:"balance"`
// }
// type Trace struct {
// 	Action              TraceAction `json:"action"`
// 	BlockHash           common.Hash `json:"blockHash"`
// 	BlockNumber         int64       `json:"blockNumber"`
// 	Result              TraceResult `json:"result"`
// 	Subtraces           int         `json:"subtraces"`
// 	TraceAddress        []int       `json:"traceAddress"`
// 	TransactionHash     common.Hash `json:"transactionHash"`
// 	TransactionPosition int         `json:"transactionPosition"`
// 	Type                string      `json:"type"`
// 	tag                 string
// }

// func (s *Trace) Tag() string {
// 	if s.tag == "" {
// 		s.tag = s.Action.CallType + "_" + gstr.JoinAny(s.TraceAddress, "_")
// 	}
// 	return s.tag
// }

// type TraceResult struct {
// 	GasUsed string `json:"gasUsed"`
// 	Output  string `json:"output`
// 	Address string `json:"address"`
// 	Code    string `json:"code"`
// }

// func (s *EthTrace) GetTraceTransfer(ctx context.Context, block *ethtypes.Block) ([]*entity.SynctransferTransfer, error) {
// 	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeOut)
// 	defer cancel()

// 	traces, err := s.traceBlock(ctx, block.Number())
// 	if err != nil {
// 		return nil, errors.New(fmt.Sprintln("getTraceBlock:", block.Number(), err))
// 	}

// 	t := s.processInTxns(ctx, block, traces)
// 	return t, nil
// }
// func (s *EthTrace) traceBlock(ctx context.Context, number *big.Int) ([]*Trace, error) {
// 	var head []*Trace
// 	err := s.cli.Client().CallContext(ctx, &head, "trace_block", toBlockNumArg(number))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return head, err
// }

// func (s *EthTrace) processInTxns(ctx context.Context, block *ethtypes.Block, traces []*Trace) []*entity.SynctransferTransfer {
// 	///

// 	////
// 	filtertrace := []*Trace{}
// 	for _, trace := range traces {
// 		if trace.Action.CallType == "call" {
// 			if trace.Action.Value.String() != "0x0" && trace.Action.Input == "0x" && len(trace.TraceAddress) > 0 {
// 				filtertrace = append(filtertrace, trace)
// 			}
// 		}
// 	}

// 	//// fill transfer
// 	transfers := []*entity.SynctransferTransfer{}
// 	for _, trace := range filtertrace {
// 		tx := block.Transaction(trace.TransactionHash)
// 		if tx == nil {
// 			g.Log().Warning(ctx, "tx is nil")
// 			continue
// 		}
// 		///notice: drop intx if traceaddress > 5555
// 		if len(trace.TraceAddress) > 8 {
// 			g.Log().Warning(ctx, "intx too long traceaddress:", trace)
// 			continue
// 		}
// 		/////
// 		transfer := &entity.SynctransferTransfer{
// 			ChainId:   s.chainId,
// 			Height:    trace.BlockNumber,
// 			BlockHash: trace.BlockHash.String(),
// 			Ts:        int64(block.Time()),
// 			TxHash:    trace.TransactionHash.String(),
// 			TxIdx:     trace.TransactionPosition,
// 			From:      trace.Action.From.String(),
// 			To:        trace.Action.To.String(),
// 			Contract:  "",
// 			Value:     trace.Action.Value.ToInt().String(),
// 			Gas:       trace.Action.Gas,
// 			GasPrice:  "0",
// 			LogIdx:    -1,
// 			Nonce:     int64(tx.Nonce()),
// 			Kind:      "external",
// 			Status:    0,
// 			Removed:   false,
// 			TraceTag:  trace.Tag(),
// 		}
// 		transfers = append(transfers, transfer)
// 	}

// 	return transfers
// }
