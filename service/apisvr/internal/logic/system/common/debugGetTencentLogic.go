package common

import (
	"context"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DebugGetTencentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDebugGetTencentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DebugGetTencentLogic {
	return &DebugGetTencentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DebugGetTencentLogic) DebugGetTencent() error {
	// todo: add your logic here and delete this line

	return nil
}
