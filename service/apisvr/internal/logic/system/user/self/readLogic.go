package self

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/user/info"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/share/ctxs"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.UserSelfReadReq) (resp *types.UserInfo, err error) {
	var uc = ctxs.GetUserCtx(l.ctx)
	return info.NewReadLogic(l.ctx, l.svcCtx).Read(&types.UserInfoReadReq{
		UserID:     uc.UserID,
		WithRoles:  req.WithRoles,
		WithTenant: req.WithTenant,
	})
}
