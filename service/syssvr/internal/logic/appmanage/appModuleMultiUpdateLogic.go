package appmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppModuleMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppModuleMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppModuleMultiUpdateLogic {
	return &AppModuleMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppModuleMultiUpdateLogic) AppModuleMultiUpdate(in *sys.AppModuleMultiUpdateReq) (*sys.Empty, error) {
	err := relationDB.NewAppModuleRepo(l.ctx).MultiUpdate(l.ctx, in.Code, in.ModuleCodes)
	return &sys.Empty{}, err
}
