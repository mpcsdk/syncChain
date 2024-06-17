package packed

import (
	"syncChain/internal/logic/chaindata"
	"syncChain/internal/logic/db"
	msg "syncChain/internal/logic/eventsender"
	"syncChain/internal/service"
)

func init() {
	service.RegisterEvnetSender(msg.NewMsg())
	service.RegisterDB(db.New())
	service.RegisterChainData(chaindata.New())
}
