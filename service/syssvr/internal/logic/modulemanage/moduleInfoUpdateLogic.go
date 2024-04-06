package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoUpdateLogic {
	return &ModuleInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoUpdateLogic) ModuleInfoUpdate(in *sys.ModuleInfo) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewModuleInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.Name = in.Name
	old.Path = in.Path
	old.Order = in.Order
	old.Url = in.Url
	old.Icon = in.Icon
	old.Body = in.Body.Value
	old.HideInMenu = in.HideInMenu
	old.Type = in.Type
	old.SubType = in.SubType
	if in.Desc != nil {
		old.Desc = in.Desc.Value
	}

	err = relationDB.NewModuleInfoRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
