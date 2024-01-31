package self

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/module/info"
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/ctxs"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModuleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleIndexLogic {
	return &ModuleIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModuleIndexLogic) ModuleIndex() (resp *types.ModuleInfoIndexResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	roleID := uc.RoleID
	if roleID == 0 {
		return nil, nil
	}
	var moduleCodes []string
	if !uc.IsAdmin {
		codes, err := l.svcCtx.RoleRpc.RoleModuleIndex(l.ctx, &role.RoleModuleIndexReq{AppCode: uc.AppCode, Id: roleID})
		if err != nil {
			return nil, err
		}
		if len(codes.ModuleCodes) == 0 {
			return nil, nil
		}
		moduleCodes = codes.ModuleCodes
	}

	ret, err := l.svcCtx.TenantRpc.TenantAppModuleIndex(l.ctx, &sys.TenantModuleIndexReq{Code: uc.TenantCode, AppCode: uc.AppCode, ModuleCodes: moduleCodes})
	if err != nil {
		return nil, err
	}

	return &types.ModuleInfoIndexResp{
		List: info.ToModuleInfosApi(ret.List),
	}, nil
}
