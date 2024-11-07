package self

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定账号
func NewBindAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindAccountLogic {
	return &BindAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindAccountLogic) BindAccount(req *types.UserBindAccountReq) error {
	_, err := l.svcCtx.UserRpc.UserBindAccount(l.ctx, &sys.UserBindAccountReq{
		Account: req.Account,
		Type:    req.Type,
		Code:    req.Code,
		CodeID:  req.CodeID,
	})
	return err
}
