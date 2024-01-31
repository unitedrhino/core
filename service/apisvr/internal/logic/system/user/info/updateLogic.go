package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/user"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UserInfo) error {
	_, err := l.svcCtx.UserRpc.UserInfoUpdate(l.ctx, user.UserInfoToRpc(req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.user.upadte failure err=%+v", utils.FuncName(), er)
		return er
	}
	return nil
}
