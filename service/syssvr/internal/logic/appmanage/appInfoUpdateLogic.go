package appmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

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

func (l *AppInfoUpdateLogic) AppInfoUpdate(in *sys.AppInfo) (*sys.Response, error) {
	err := relationDB.NewAppInfoRepo(l.ctx).Update(l.ctx, ToAppInfoPo(in))
	return &sys.Response{}, err
}
