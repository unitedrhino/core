package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuCreateLogic {
	return &ModuleMenuCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuCreateLogic) ModuleMenuCreate(in *sys.MenuInfo) (*sys.WithID, error) {
	if err := CheckModule(l.ctx, in.ModuleCode); err != nil {
		return nil, err
	}
	if in.Type == 0 {
		in.Type = 1
	}
	if in.ParentID == 0 {
		in.ParentID = 1
	}
	if in.Order == 0 {
		in.Order = 1
	}
	if in.HideInMenu == 0 {
		in.HideInMenu = 1
	}
	po := logic.ToMenuInfoPo(in)
	relationDB.NewMenuInfoRepo(l.ctx).Insert(l.ctx, po)
	return &sys.WithID{Id: po.ID}, nil
}
