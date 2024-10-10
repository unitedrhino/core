package appmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppModuleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppModuleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppModuleIndexLogic {
	return &AppModuleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppModuleIndexLogic) AppModuleIndex(in *sys.AppModuleIndexReq) (*sys.AppModuleIndexResp, error) {
	list, err := relationDB.NewAppModuleRepo(l.ctx).FindByFilter(l.ctx, relationDB.AppModuleFilter{AppCodes: []string{in.Code}}, nil)
	if err != nil {
		return nil, err
	}
	var moduleCodes []string
	for _, v := range list {
		moduleCodes = append(moduleCodes, v.ModuleCode)
	}
	return &sys.AppModuleIndexResp{ModuleCodes: moduleCodes}, nil
}
