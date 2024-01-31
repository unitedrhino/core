package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ModuleInfo) (resp *types.WithID, err error) {
	ret, err := l.svcCtx.ModuleRpc.ModuleInfoCreate(l.ctx, ToModuleInfoRpc(req))
	return logic.SysToWithIDTypes(ret), err
}
