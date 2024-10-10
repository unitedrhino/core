package accessmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoIndexLogic {
	return &ApiInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiInfoIndexLogic) ApiInfoIndex(in *sys.ApiInfoIndexReq) (*sys.ApiInfoIndexResp, error) {
	f := relationDB.ApiInfoFilter{
		ApiIDs:     in.ApiIDs,
		Route:      in.Route,
		Method:     in.Method,
		Name:       in.Name,
		AccessCode: in.AccessCode,
		AuthType:   in.AuthType,
	}
	pos, err := relationDB.NewApiInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewApiInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	var ais []*sys.ApiInfo
	for _, v := range pos {
		ais = append(ais, ToApiInfoPb(v))
	}
	return &sys.ApiInfoIndexResp{
		List:  ais,
		Total: total,
	}, nil
}
