package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system"
	"gitee.com/i-Things/core/shared/ctxs"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.WithIDOrCode) (resp *types.TenantInfo, err error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.TenantRpc.TenantInfoRead(l.ctx, system.ToSysWithIDCode(req))
	return system.ToTenantInfoTypes(ret), err
}
