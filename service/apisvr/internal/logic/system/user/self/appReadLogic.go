package self

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppReadLogic {
	return &AppReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppReadLogic) AppRead(req *types.UserSelfAppReadReq) (resp *types.UserSelfAppReadResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppIndex(ctxs.WithRoot(l.ctx), &sys.TenantAppIndexReq{
		Type:    req.Type,
		SubType: req.SubType,
		AppID:   req.AppID,
	})
	if err != nil {
		return nil, err
	}
	if len(ret.List) == 0 {
		return &types.UserSelfAppReadResp{}, nil
	}
	appInfo, err := l.svcCtx.AppRpc.AppInfoRead(l.ctx, &sys.WithIDCode{Code: ret.List[0].AppCode})
	if err != nil {
		return nil, err
	}
	resp = &types.UserSelfAppReadResp{Code: appInfo.Code, Name: appInfo.Name}
	for _, v := range ret.List {
		resp.TenantCodes = append(resp.TenantCodes, v.Code)
	}
	return
}
