package main

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/timedjobdirect"
	"gitee.com/unitedrhino/share/utils"
)

func main() {
	defer utils.Recover(context.Background())
	ctx := timedjobdirect.GetSvcCtx()
	timedjobdirect.Run(ctx)
}
