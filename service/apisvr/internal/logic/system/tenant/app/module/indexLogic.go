package module

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/module/info"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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

func (l *IndexLogic) Index(req *types.TenantModuleIndexReq) (resp *types.TenantModuleIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppModuleIndex(l.ctx, &sys.TenantModuleIndexReq{
		Page:    logic.ToSysPageRpc(req.Page),
		Code:    req.Code,
		AppCode: req.AppCode,
	})
	if err != nil {
		return nil, err
	}
	return &types.TenantModuleIndexResp{
		List: info.ToModuleInfosApi(ret.List),
	}, nil
}
