package accessmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/sysExport"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoDeleteLogic {
	return &ApiInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiInfoDeleteLogic) ApiInfoDelete(in *sys.WithID) (*sys.Empty, error) {
	db := relationDB.NewApiInfoRepo(l.ctx)
	old, err := db.FindOne(l.ctx, in.Id)
	if err != nil {
		return &sys.Empty{}, err
	}
	err = relationDB.NewApiInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	if err == nil {
		l.svcCtx.ApiCache.SetData(l.ctx, sysExport.GenApiCacheKey(old.Method, old.Route), nil)
	}
	return &sys.Empty{}, err
}
