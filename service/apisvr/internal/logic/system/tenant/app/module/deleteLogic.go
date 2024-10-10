package module

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

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

func (l *DeleteLogic) Delete(req *types.TenantModuleWithIDOrCode) error {
	_, err := l.svcCtx.TenantRpc.TenantAppModuleDelete(l.ctx, &sys.TenantModuleWithIDOrCode{
		Code:       req.Code,
		AppCode:    req.AppCode,
		ModuleCode: req.ModuleCode,
		Id:         req.ID,
	})
	return err
}
