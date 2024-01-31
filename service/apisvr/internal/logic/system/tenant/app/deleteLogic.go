package app

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/ctxs"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

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

func (l *DeleteLogic) Delete(req *types.TenantAppWithIDOrCode) error {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return err
	}
	_, err := l.svcCtx.TenantRpc.TenantAppDelete(l.ctx, &sys.TenantAppWithIDOrCode{
		Code:    req.Code,
		AppCode: req.AppCode,
		Id:      req.ID,
	})
	return err
}
