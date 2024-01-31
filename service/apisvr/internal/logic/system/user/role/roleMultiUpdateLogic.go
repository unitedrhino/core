package role

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/errors"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMultiUpdateLogic {
	return &RoleMultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleMultiUpdateLogic) RoleMultiUpdate(req *types.UserRoleMultiUpdateReq) error {
	uc := ctxs.GetUserCtx(l.ctx)
	//这里需要判断是否是租户下的超级管理员,只有租户下的超级管理员才能修改角色
	ti, err := l.svcCtx.TenantRpc.TenantInfoRead(l.ctx, &sys.WithIDCode{Code: uc.TenantCode})
	if err != nil {
		return err
	}
	if ti.AdminUserID != uc.UserID {
		return errors.Permissions.AddDetail("非超级管理员不能修改角色")
	}
	_, err = l.svcCtx.UserRpc.UserRoleMultiUpdate(l.ctx, &sys.UserRoleMultiUpdateReq{UserID: req.UserID, RoleIDs: req.RoleIDs})
	return err
}
