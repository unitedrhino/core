package timedschedulerlogic

import (
	"context"

	"gitee.com/i-Things/core/service/timed/timedschedulersvr/internal/svc"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *timedscheduler.Empty) (*timedscheduler.Empty, error) {
	// todo: add your logic here and delete this line

	return &timedscheduler.Empty{}, nil
}
