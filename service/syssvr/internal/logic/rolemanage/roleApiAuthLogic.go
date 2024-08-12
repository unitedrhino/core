package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/domain/access"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
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
	api, err := relationDB.NewApiInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ApiInfoFilter{
		Route:  in.Path,
		Method: in.Method,
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Permissions
		}
		return nil, err
	}
	if api.Access == nil {
		return nil, errors.Permissions
	}
	if api.Access.AuthType == access.AuthTypeSupper && uc.TenantCode != def.TenantCodeDefault {
		return nil, errors.Permissions.AddDetail("只有租户管理员才能操作")
	}
	_, err = relationDB.NewRoleAccessRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.RoleAccessFilter{RoleIDs: uc.RoleIDs, AccessCodes: []string{api.AccessCode}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Permissions
		}
		return nil, err
	}
	return &sys.RoleApiAuthResp{BusinessType: api.BusinessType, Name: api.Name}, errors.Permissions.AddMsg("无权限")
}
