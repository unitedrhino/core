package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
	ro, err := l.RiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Error("RoleInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}
	if in.Name == "" || ro.Name == "超级管理员" {
		in.Name = ro.Name
	}

	if in.Desc == "" {
		in.Desc = ro.Desc
	}

	if in.Status == 0 {
		in.Status = ro.Status
	}

	err = l.RiDB.Update(l.ctx, &relationDB.SysRoleInfo{
		ID:     in.Id,
		Name:   in.Name,
		Desc:   in.Desc,
		Status: in.Status,
	})
	if err != nil {
		return nil, err
	}
	return &sys.Empty{}, nil
}
