package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoDeleteLogic {
	return &ModuleInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoDeleteLogic) ModuleInfoDelete(in *sys.WithIDCode) (*sys.Response, error) {
	f := relationDB.ModuleInfoFilter{ID: in.Id}
	if in.Code != "" {
		f.Codes = []string{in.Code}
	}
	err := relationDB.NewModuleInfoRepo(l.ctx).DeleteByFilter(l.ctx, f)
	return &sys.Response{}, err
}
