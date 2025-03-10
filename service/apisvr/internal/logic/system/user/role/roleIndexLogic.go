package role

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/role"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleIndexLogic {
	return &RoleIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleIndexLogic) RoleIndex(req *types.UserRoleIndexReq) (resp *types.UserRoleIndexResp, err error) {
	ret, err := l.svcCtx.UserRpc.UserRoleIndex(l.ctx, &sys.UserRoleIndexReq{UserID: req.UserID})
	if err != nil {
		return nil, err
	}

	return &types.UserRoleIndexResp{
		List:     role.ToRoleInfosTypes(ret.List),
		PageResp: logic.ToPageResp(nil, ret.Total),
	}, nil
}
