package menu

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.TenantAppMenuIndexReq) (resp *types.TenantAppMenuIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppMenuIndex(l.ctx, &sys.TenantAppMenuIndexReq{
		AppCode:    req.AppCode,
		Code:       req.Code,
		ModuleCode: req.ModuleCode,
		IsRetTree:  req.IsRetTree,
	})

	return &types.TenantAppMenuIndexResp{
		List: system.ToTenantAppMenusApi(ret.List),
	}, nil
}
