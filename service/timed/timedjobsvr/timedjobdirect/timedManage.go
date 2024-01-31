package timedjobdirect

import (
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	server "gitee.com/i-Things/core/service/timed/timedjobsvr/internal/server/timedmanage"
)

var (
	jobSvr timedmanage.TimedManage
)

func NewTimedJob(runSvr bool) timedmanage.TimedManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	svr := timedmanage.NewDirectTimedManage(svcCtx, server.NewTimedManageServer(svcCtx))
	return svr
}
