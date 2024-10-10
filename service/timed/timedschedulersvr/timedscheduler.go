package main

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/timedschedulerdirect"
	"gitee.com/unitedrhino/share/utils"
)

func main() {
	defer utils.Recover(context.Background())
	ctx := timedschedulerdirect.GetSvcCtx()
	timedschedulerdirect.Run(ctx)
}
