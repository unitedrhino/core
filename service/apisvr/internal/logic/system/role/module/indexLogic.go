package module

import (
	"context"
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

func (l *IndexLogic) Index(req *types.RoleModuleIndexReq) (resp *types.RoleModuleIndexResp, err error) {
	ret, err := l.svcCtx.RoleRpc.RoleModuleIndex(l.ctx, &sys.RoleModuleIndexReq{Id: req.ID, AppCode: req.AppCode})
	if err != nil {
		return nil, err
	}
	return &types.RoleModuleIndexResp{
		ModuleCodes: ret.ModuleCodes,
	}, nil
}
