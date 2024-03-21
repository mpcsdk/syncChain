package main

import (
	_ "syncChain/internal/packed"

	_ "syncChain/internal/logic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"syncChain/internal/cmd"
)

var SyncBlock bool = false

func main() {

	///

	g.Log().Async(true)
	cmd.Main.Run(gctx.GetInitCtx())
}
