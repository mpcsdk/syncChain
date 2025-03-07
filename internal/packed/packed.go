package packed

import (
	"context"
	"os"
	"syncChain/internal/logic/chaindata"
	"syncChain/internal/logic/db"
	msg "syncChain/internal/logic/eventsender"
	"syncChain/internal/service"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gproc"
)

func signalHandlerForExit(sig os.Signal) {
	g.Log().Info(context.Background(), "shutting down due to signal:", sig.String())
	service.ChainData().Stop()
	for {
		if service.ChainData().IsRunning() {
			g.Log().Info(context.Background(), "wait exit... ")
			time.Sleep(time.Second)
		} else {
			return
		}
	}
}
func init() {
	service.RegisterEvnetSender(msg.NewMsg())
	service.RegisterDB(db.New())
	service.RegisterChainData(chaindata.New())
	gproc.AddSigHandlerShutdown(
		signalHandlerForExit,
	)
}
