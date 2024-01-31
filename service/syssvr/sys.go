// 系统管理模块-syssvr
package main

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/sysdirect"
	"gitee.com/i-Things/share/utils"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := sysdirect.GetSvcCtx()
	sysdirect.Run(svcCtx)
}
