package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAppMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAppMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAppMultiUpdateLogic {
	return &RoleAppMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAppMultiUpdateLogic) RoleAppMultiUpdate(in *sys.RoleAppMultiUpdateReq) (*sys.Response, error) {
	err := relationDB.NewRoleAppRepo(l.ctx).MultiUpdate(l.ctx, in.Id, in.AppCodes)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
