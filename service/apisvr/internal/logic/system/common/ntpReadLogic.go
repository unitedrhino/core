package common

import (
	"context"
	"time"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NtpReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNtpReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NtpReadLogic {
	return &NtpReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NtpReadLogic) NtpRead(req *types.NtpReadReq) (resp *types.NtpReadResp, err error) {
	resp = &types.NtpReadResp{
		DeviceSendTime: req.DeviceSendTime,
		ServerSendTime: time.Now().UnixMilli(),
		ServerRecvTime: time.Now().UnixMilli(),
	}
	return
}
