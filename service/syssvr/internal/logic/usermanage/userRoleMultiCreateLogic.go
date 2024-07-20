package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleMultiCreateLogic {
	return &UserRoleMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleMultiCreateLogic) UserRoleMultiCreate(in *sys.UserRoleMultiUpdateReq) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	if len(in.RoleCodes) != 0 {
		rs, err := relationDB.NewRoleInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.RoleInfoFilter{
			Codes: in.RoleCodes,
		}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range rs {
			in.RoleIDs = append(in.RoleIDs, v.ID)
		}
	}
	var datas []*relationDB.SysUserRole
	for _, v := range in.RoleIDs {
		datas = append(datas, &relationDB.SysUserRole{
			RoleID: v,
			UserID: in.UserID,
		})
	}
	err := relationDB.NewUserRoleRepo(l.ctx).MultiInsert(l.ctx, datas)
	return &sys.Empty{}, err
}
