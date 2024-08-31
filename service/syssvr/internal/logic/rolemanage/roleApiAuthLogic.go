package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/domain/access"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/service/syssvr/sysExport"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleApiAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiAuthLogic {
	return &RoleApiAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleApiAuthLogic) RoleApiAuth(in *sys.RoleApiAuthReq) (*sys.RoleApiAuthResp, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	//if uc.IsAdmin {
	//	return &sys.RoleApiAuthResp{}, nil
	//}
	//api, err := relationDB.NewApiInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ApiInfoFilter{
	//	Route:  in.Path,
	//	Method: in.Method,
	//})
	api, err := l.svcCtx.ApiCache.GetData(l.ctx, sysExport.GenApiCacheKey(in.Method, in.Path))
	if err != nil {
		if errors.Cmp(err, errors.NotFind) { //没有导入的不校验
			return &sys.RoleApiAuthResp{}, nil
		}
		return nil, err
	}
	if uc.IsSuperAdmin {
		return &sys.RoleApiAuthResp{BusinessType: api.BusinessType, Name: api.Name}, nil
	}
	if api.Access == nil {
		return nil, errors.Permissions
	}
	if api.Access.IsNeedAuth != def.True {
		return &sys.RoleApiAuthResp{BusinessType: api.BusinessType, Name: api.Name}, nil
	}
	if api.Access.AuthType == access.AuthTypeSupper && uc.TenantCode != def.TenantCodeDefault {
		return nil, errors.Permissions.AddDetail("只有租户管理员才能操作")
	}
	ras, err := l.svcCtx.RoleAccessCache.GetData(l.ctx, api.AccessCode)
	if err != nil {
		return nil, err
	}
	for _, roleID := range uc.RoleIDs {
		if _, ok := (*ras)[roleID]; ok {
			return &sys.RoleApiAuthResp{BusinessType: api.BusinessType, Name: api.Name}, nil
		}
	}
	return nil, errors.Permissions
}
