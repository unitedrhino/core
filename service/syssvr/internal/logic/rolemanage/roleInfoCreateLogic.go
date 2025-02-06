package rolemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
}

func NewRoleInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleInfoCreateLogic {
	return &RoleInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func (l *RoleInfoCreateLogic) RoleInfoCreate(in *sys.RoleInfo) (*sys.WithID, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	po := relationDB.SysRoleInfo{
		Name:   in.Name,
		Desc:   in.Desc.GetValue(),
		Status: in.Status,
		Code:   in.Code,
	}
	if in.Code == "" {
		in.Code = utils.Random(10, 1)
	}
	err := l.RiDB.Insert(l.ctx, &po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}
