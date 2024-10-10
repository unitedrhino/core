// 系统管理模块-syssvr
package main

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/sysdirect"
	"gitee.com/unitedrhino/share/utils"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := sysdirect.GetSvcCtx()
	sysdirect.Run(svcCtx)
}
