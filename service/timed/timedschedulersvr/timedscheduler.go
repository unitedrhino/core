package main

import (
	"context"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/timedschedulerdirect"
	"gitee.com/i-Things/core/shared/utils"
)

func main() {
	defer utils.Recover(context.Background())
	ctx := timedschedulerdirect.GetSvcCtx()
	timedschedulerdirect.Run(ctx)
}
