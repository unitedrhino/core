package access

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取开放认证详情
func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.WithID) (resp *types.DataOpenAccess, err error) {
	ret, err := l.svcCtx.DataM.DataOpenAccessRead(l.ctx, utils.Copy[sys.WithID](req))
	if err != nil {
		return nil, err
	}
	resp = utils.Copy[types.DataOpenAccess](ret)
	u, err := l.svcCtx.UserCache.GetData(l.ctx, ret.UserID)
	if err != nil {
		return
	}
	resp.User = utils.Copy[types.UserCore](u)
	return utils.Copy[types.DataOpenAccess](ret), nil
}
