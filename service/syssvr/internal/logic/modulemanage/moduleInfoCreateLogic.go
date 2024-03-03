package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoCreateLogic {
	return &ModuleInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoCreateLogic) ModuleInfoCreate(in *sys.ModuleInfo) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	if in.Type == 0 {
		in.Type = 1
	}
	if in.Order == 0 {
		in.Order = 1
	}
	if in.HideInMenu == 0 {
		in.HideInMenu = 1
	}
	po := logic.ToModuleInfoPo(in)
	relationDB.NewModuleInfoRepo(l.ctx).Insert(l.ctx, po)
	return &sys.WithID{Id: po.ID}, nil
}
