package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.UserInfoDeleteReq) error {
	_, err := l.svcCtx.UserRpc.UserInfoDelete(l.ctx, &sys.UserInfoDeleteReq{
		UserID: req.UserID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.InfoDelete req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
