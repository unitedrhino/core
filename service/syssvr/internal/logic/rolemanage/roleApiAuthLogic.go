package rolemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
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

func (l *RoleApiAuthLogic) RoleApiAuth(in *sys.RoleApiAuthReq) (*sys.Empty, error) {
	return &sys.Empty{}, nil

	//todo 待实现
	//uc := ctxs.GetUserCtx(l.ctx)
	//api, err := relationDB.NewRoleApiRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.RoleApiFilter{
	//	AppCode:   uc.AppCode,
	//	Route:     in.Path,
	//	Method:    in.Method,
	//	WithRoles: true,
	//})
	//if err != nil && !errors.Cmp(err, errors.NotFind) {
	//	return nil, err
	//}
	//
	//if errors.Cmp(err, errors.NotFind) { //如果没有找到,可能是不需要鉴权的接口,需要去模块接口里查询
	//	api, err := relationDB.NewApiInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ApiInfoFilter{
	//		Route:  in.Path,
	//		Method: in.Method,
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//	if api.IsNeedAuth == def.True { //需要鉴权的接口,但是没有查询到,说明这个租户没有权限
	//		return nil, errors.Permissions.AddDetail("权限不足")
	//	}
	//	return &sys.Empty{}, nil
	//}
	//if uc.IsAdmin { //如果是租户管理员,则有权限
	//	return &sys.Empty{}, nil
	//}
	//for _, v := range api.Roles {
	//	if v.RoleID == in.RoleID {
	//		return &sys.Empty{}, nil
	//	}
	//}
	return nil, errors.Permissions.AddDetail("权限不足")
}
