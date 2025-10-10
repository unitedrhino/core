package modulemanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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
	if in.IsProject != 0 {
		old.IsProject = in.IsProject
	}
	if in.IsPlatform != 0 {
		old.IsPlatform = in.IsPlatform
	}
	if in.Desc != nil {
		old.Desc = in.Desc.Value
	}
	if in.HomeMenuID == 1 {
		old.HomeMenuID = in.HomeMenuID
	}
	if in.HomeMenuID > 1 {
		_, err := relationDB.NewMenuInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.MenuInfoFilter{
			ModuleCode: old.Code, MenuID: in.HomeMenuID,
		})
		if err != nil {
			return nil, errors.Parameter.AddMsg("请选择模块下有的菜单")
		}
		old.HomeMenuID = in.HomeMenuID
	}

	err = relationDB.NewModuleInfoRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
