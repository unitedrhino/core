package usermanagelogic

import (
	"context"
	rolemanagelogic "gitee.com/i-Things/core/service/syssvr/internal/logic/rolemanage"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleIndexLogic {
	return &UserRoleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleIndexLogic) UserRoleIndex(in *sys.UserRoleIndexReq) (*sys.UserRoleIndexResp, error) {
	ur, err := relationDB.NewUserRoleRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserRoleFilter{UserID: in.UserID}, nil)
	if err != nil {
		return nil, err
	}
	var roleIDs []int64
	for _, v := range ur {
		roleIDs = append(roleIDs, v.RoleID)
	}
	rs, err := relationDB.NewRoleInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.RoleInfoFilter{IDs: roleIDs}, nil)
	if err != nil {
		return nil, err
	}
	return &sys.UserRoleIndexResp{List: rolemanagelogic.ToRoleInfosRpc(rs), Total: int64(len(roleIDs))}, nil
}
