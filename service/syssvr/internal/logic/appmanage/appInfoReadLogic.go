package appmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoReadLogic {
	return &AppInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoReadLogic) AppInfoRead(in *sys.WithIDCode) (*sys.AppInfo, error) {
	f := relationDB.AppInfoFilter{ID: in.Id}
	if in.Code != "" {
		f.Codes = []string{in.Code}
	}
	ret, err := relationDB.NewAppInfoRepo(l.ctx).FindOneByFilter(l.ctx, f)
	return ToAppInfoPb(ret), err
}
