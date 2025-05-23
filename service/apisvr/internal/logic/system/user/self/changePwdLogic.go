package self

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePwdLogic {
	return &ChangePwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePwdLogic) ChangePwd(req *types.UserChangePwdReq) error {
	_, err := l.svcCtx.UserRpc.UserChangePwd(l.ctx, &sys.UserChangePwdReq{
		Type:        req.Type,
		Password:    req.Password,
		OldPassword: req.OldPassword,
		Code:        req.Code,
		CodeID:      req.CodeID,
	})
	return err
}
