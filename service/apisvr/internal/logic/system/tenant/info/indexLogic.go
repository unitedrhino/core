package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.TenantInfoIndexReq) (resp *types.TenantInfoIndexResp, err error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.TenantRpc.TenantInfoIndex(l.ctx, &sys.TenantInfoIndexReq{
		Name: req.Name,
		Page: logic.ToSysPageRpc(req.Page),
		Code: req.Code,
	})
	if err != nil {
		return nil, err
	}
	var userMap = map[int64]*sys.UserInfo{}
	if req.WithAdminUser && len(ret.List) > 0 {
		var userIDs []int64
		for _, v := range ret.List {
			userIDs = append(userIDs, v.AdminUserID)
		}
		users, err := l.svcCtx.UserRpc.UserInfoIndex(ctxs.WithRoot(l.ctx), &sys.UserInfoIndexReq{
			UserIDs: userIDs,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range users.List {
			userMap[v.UserID] = v
		}
	}
	return &types.TenantInfoIndexResp{
		Total: ret.Total,
		List:  system.ToTenantInfosTypes(ret.List, userMap),
	}, nil
}
