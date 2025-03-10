package core

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

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
		Name:    req.Name,
		Page:    logic.ToSysPageRpc(req.Page),
		Code:    req.Code,
		AppCode: req.AppCode,
		Status:  def.True,
	})
	if err != nil {
		return nil, err
	}
	return &types.TenantCoreIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     system.ToTenantCoresTypes(ret.List),
	}, nil
}
