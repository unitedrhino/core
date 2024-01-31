package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleModuleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleModuleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleModuleIndexLogic {
	return &RoleModuleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleModuleIndexLogic) RoleModuleIndex(in *sys.RoleModuleIndexReq) (*sys.RoleModuleIndexResp, error) {
	ret, err := relationDB.NewRoleModuleRepo(l.ctx).FindByFilter(l.ctx, relationDB.RoleModuleFilter{
		RoleIDs: []int64{in.Id},
		AppCode: in.AppCode,
	}, nil)
	if err != nil {
		return nil, err
	}
	var moduleCodes []string
	if len(ret) != 0 {
		for _, v := range ret {
			moduleCodes = append(moduleCodes, v.ModuleCode)
		}
	}
	return &sys.RoleModuleIndexResp{ModuleCodes: moduleCodes}, nil
}
