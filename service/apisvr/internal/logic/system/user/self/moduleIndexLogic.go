package self

import (
	"context"
	role "gitee.com/unitedrhino/core/service/syssvr/client/rolemanage"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

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

func (l *ModuleIndexLogic) ModuleIndex(req *types.UserModuleIndexReq) (resp *types.TenantModuleInfoIndexResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	var moduleCodes []string
	if !uc.IsSuperAdmin {
		codes, err := l.svcCtx.RoleRpc.RoleModuleIndex(l.ctx, &role.RoleModuleIndexReq{AppCode: uc.AppCode, Ids: uc.RoleIDs})
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
	resp = &types.TenantModuleInfoIndexResp{
		List: utils.CopySlice[types.TenantModuleInfo](ret.List),
	}
	if req.WithMenus {
		for _, m := range resp.List {
			r, err := NewMenuIndexLogic(l.ctx, l.svcCtx).MenuIndex(&types.UserMenuIndexReq{ModuleCode: m.Code})
			if err != nil {
				return nil, err
			}
			m.Menus = r.List
		}
	}
	return resp, nil
}
