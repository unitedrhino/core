package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleModuleMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleModuleMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleModuleMultiUpdateLogic {
	return &RoleModuleMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleModuleMultiUpdateLogic) RoleModuleMultiUpdate(in *sys.RoleModuleMultiUpdateReq) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	err := relationDB.NewRoleModuleRepo(l.ctx).MultiUpdate(l.ctx, in.Id, in.AppCode, in.ModuleCodes)
	return &sys.Empty{}, err
}
