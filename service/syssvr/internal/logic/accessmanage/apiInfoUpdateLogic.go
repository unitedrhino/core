package accessmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoUpdateLogic {
	return &ApiInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiInfoUpdateLogic) ApiInfoUpdate(in *sys.ApiInfo) (*sys.Empty, error) {
	old, err := relationDB.NewApiInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.AccessCode = in.AccessCode
	old.Method = in.Method
	old.Route = in.Route
	old.Name = in.Name
	//old.BusinessType = in.BusinessType
	old.RecordLogMode = in.RecordLogMode
	old.Desc = in.Desc
	//old.AuthType = in.AuthType
	err = relationDB.NewApiInfoRepo(l.ctx).Update(l.ctx, old)
	if err == nil {
		l.svcCtx.ApiCache.SetData(l.ctx, sysExport.GenApiCacheKey(old.Method, old.Route), nil)
	}
	return &sys.Empty{}, err
}
