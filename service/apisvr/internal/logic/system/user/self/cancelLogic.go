package self

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLogic {
	return &CancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelLogic) Cancel() error {
	uc := ctxs.GetUserCtx(l.ctx)
	_, err := l.svcCtx.UserRpc.UserInfoDelete(l.ctx, &sys.UserInfoDeleteReq{
		UserID: uc.UserID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.InfoDelete err=%+v", utils.FuncName(), er)
		return er
	}
	return nil
}
