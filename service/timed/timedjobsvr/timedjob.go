package main

import (
	"context"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/timedjobdirect"
	"gitee.com/i-Things/share/utils"
)

func main() {
	defer utils.Recover(context.Background())
	ctx := timedjobdirect.GetSvcCtx()
	timedjobdirect.Run(ctx)
}
