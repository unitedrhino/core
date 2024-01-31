package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
}

func NewRoleInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleInfoIndexLogic {
	return &RoleInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func (l *RoleInfoIndexLogic) RoleInfoIndex(in *sys.RoleInfoIndexReq) (*sys.RoleInfoIndexResp, error) {
	f := relationDB.RoleInfoFilter{
		Name:   in.Name,
		Status: in.Status,
	}
	ros, err := l.RiDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := l.RiDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return &sys.RoleInfoIndexResp{
		List:  ToRoleInfosRpc(ros),
		Total: total,
	}, nil
}
