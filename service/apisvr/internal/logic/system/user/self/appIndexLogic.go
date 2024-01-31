package self

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/app/info"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppIndexLogic {
	return &AppIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppIndexLogic) AppIndex() (resp *types.AppInfoIndexResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	roleID := uc.RoleID
	if roleID == 0 {
		return nil, nil
	}
	var appCodes []string
	if !uc.IsAdmin {
		as, err := l.svcCtx.RoleRpc.RoleAppIndex(l.ctx, &sys.RoleAppIndexReq{Id: roleID})
		if err != nil {
			return nil, err
		}
		appCodes = as.AppCodes
		if len(appCodes) == 0 {
			return nil, nil
		}
	}

	ret, err := l.svcCtx.TenantRpc.TenantAppIndex(l.ctx, &sys.TenantAppIndexReq{Code: uc.TenantCode, AppCodes: appCodes})

	return &types.AppInfoIndexResp{
		List: info.ToAppInfosTypes(ret.List),
	}, nil
}
