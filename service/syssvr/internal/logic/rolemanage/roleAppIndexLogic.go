package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAppIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAppIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAppIndexLogic {
	return &RoleAppIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAppIndexLogic) RoleAppIndex(in *sys.RoleAppIndexReq) (*sys.RoleAppIndexResp, error) {
	var appCodes []string
	ra, err := relationDB.NewRoleAppRepo(l.ctx).FindByFilter(l.ctx, relationDB.RoleAppFilter{RoleID: in.Id}, nil)
	if err != nil {
		return nil, err
	}
	for _, v := range ra {
		appCodes = append(appCodes, v.AppCode)
	}
	return &sys.RoleAppIndexResp{AppCodes: appCodes, Total: int64(len(appCodes))}, nil
}
