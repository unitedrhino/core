package app

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.TenantAppIndexReq) (resp *types.TenantAppIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppIndex(l.ctx, &sys.TenantAppIndexReq{Code: req.Code})
	if err != nil {
		return nil, err
	}
	if len(ret.List) == 0 {
		return &types.TenantAppIndexResp{}, nil
	}
	appCodes := make([]string, 0)
	codeIDMap := make(map[string]int64)
	for _, v := range ret.List {
		appCodes = append(appCodes, v.AppCode)
		codeIDMap[v.AppCode] = v.Id
	}
	apps, err := l.svcCtx.AppRpc.AppInfoIndex(l.ctx, &sys.AppInfoIndexReq{
		Codes: appCodes,
	})
	if err != nil {
		return nil, err
	}
	for _, v := range apps.List {
		v.Id = codeIDMap[v.Code] //修正为关联的id
	}
	return &types.TenantAppIndexResp{
		List: utils.CopySlice[types.AppInfo](apps.List),
	}, nil
}
