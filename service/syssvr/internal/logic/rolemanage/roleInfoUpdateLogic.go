package rolemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
}

func NewRoleInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleInfoUpdateLogic {
	return &RoleInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func (l *RoleInfoUpdateLogic) RoleInfoUpdate(in *sys.RoleInfo) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	ro, err := l.RiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Error("RoleInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}

	if in.Name != "" {
		ro.Name = in.Name
	}
	if in.Desc != nil {
		ro.Desc = in.Desc.Value
	}
	if in.Status != 0 {
		ro.Status = in.Status
	}

	err = l.RiDB.Update(l.ctx, ro)
	if err != nil {
		return nil, err
	}
	return &sys.Empty{}, nil
}
