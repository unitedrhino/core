package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAccessIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAccessIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAccessIndexLogic {
	return &RoleAccessIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAccessIndexLogic) RoleAccessIndex(in *sys.RoleAccessIndexReq) (*sys.RoleAccessIndexResp, error) {
	ms, err := relationDB.NewRoleAccessRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.RoleAccessFilter{RoleIDs: []int64{in.Id}}, nil)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, v := range ms {
		ids = append(ids, v.AccessCode)
	}
	return &sys.RoleAccessIndexResp{AccessCodes: ids}, nil
}
