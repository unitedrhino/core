package self

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

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
	var uc = ctxs.GetUserCtx(l.ctx)
	req.UserID = uc.UserID
	_, err := l.svcCtx.UserRpc.UserInfoUpdate(l.ctx, &sys.UserInfoUpdateReq{Info: user.UserInfoToRpc(req)})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.user.upadte failure err=%+v", utils.FuncName(), er)
		return er
	}
	return nil
}
