package appmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoDeleteLogic {
	return &AppInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoDeleteLogic) AppInfoDelete(in *sys.WithIDCode) (*sys.Response, error) {
	f := relationDB.AppInfoFilter{ID: in.Id}
	if in.Code != "" {
		f.Codes = []string{in.Code}
	}
	err := relationDB.NewAppInfoRepo(l.ctx).DeleteByFilter(l.ctx, f)
	return &sys.Response{}, err
}
