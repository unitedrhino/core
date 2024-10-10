package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

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

func (l *ReadLogic) Read(req *types.WithIDOrCode) (resp *types.ModuleInfo, err error) {
	ret, err := l.svcCtx.ModuleRpc.ModuleInfoRead(l.ctx, system.ToSysWithIDCode(req))

	return ToModuleInfoApi(ret), err
}
