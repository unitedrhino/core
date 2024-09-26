package core

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索租户信息
func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.TenantCoreIndexReq) (resp *types.TenantCoreIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantInfoIndex(ctxs.WithRoot(l.ctx), &sys.TenantInfoIndexReq{
		Name:   req.Name,
		Page:   logic.ToSysPageRpc(req.Page),
		Code:   req.Code,
		Status: def.True,
	})
	if err != nil {
		return nil, err
	}
	return &types.TenantCoreIndexResp{
		Total: ret.Total,
		List:  system.ToTenantCoresTypes(ret.List),
	}, nil
}
