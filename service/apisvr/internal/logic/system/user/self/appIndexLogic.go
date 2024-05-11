package self

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"

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
	if len(ret.List) == 0 {
		return &types.AppInfoIndexResp{}, nil
	}
	codeIDMap := make(map[string]int64)
	for _, v := range ret.List {
		appCodes = append(appCodes, v.AppCode)
		codeIDMap[v.AppCode] = v.Id
	}
	apps, err := l.svcCtx.AppRpc.AppInfoIndex(l.ctx, &sys.AppInfoIndexReq{
		Codes: appCodes,
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	for _, v := range apps.List {
		v.Id = codeIDMap[v.Code] //修正为关联的id
	}
	return &types.AppInfoIndexResp{
		List: utils.CopySlice[types.AppInfo](apps.List),
	}, nil
}
