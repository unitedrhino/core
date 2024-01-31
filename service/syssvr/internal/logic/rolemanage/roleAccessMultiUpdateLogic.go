package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAccessMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAccessMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAccessMultiUpdateLogic {
	return &RoleAccessMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAccessMultiUpdateLogic) RoleAccessMultiUpdate(in *sys.RoleAccessMultiUpdateReq) (*sys.Response, error) {
	err := relationDB.NewRoleAccessRepo(l.ctx).MultiUpdate(l.ctx, in.Id, in.AccessCodes)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
