package opslogic

import (
	"context"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsFeedbackIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpsFeedbackIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsFeedbackIndexLogic {
	return &OpsFeedbackIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OpsFeedbackIndexLogic) OpsFeedbackIndex(in *sys.OpsFeedbackIndexReq) (*sys.OpsFeedbackIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.OpsFeedbackIndexResp{}, nil
}
