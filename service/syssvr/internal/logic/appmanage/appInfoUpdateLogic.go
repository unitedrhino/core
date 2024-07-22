package appmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoUpdateLogic {
	return &AppInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoUpdateLogic) AppInfoUpdate(in *sys.AppInfo) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewAppInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Desc != nil {
		old.Desc = in.Desc.GetValue()
	}
	if in.Type != "" {
		old.Type = in.Type
	}
	if in.SubType != "" {
		old.SubType = in.SubType
	}
	if in.MiniWx != nil {
		old.MiniWx = utils.Copy[relationDB.SysTenantThird](in.MiniWx)
	}
	err = relationDB.NewAppInfoRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
