package menu

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/module/menu"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

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

func (l *UpdateLogic) Update(req *types.TenantAppMenu) (resp *types.WithID, err error) {
	_, err = l.svcCtx.TenantRpc.TenantAppMenuUpdate(l.ctx, &sys.TenantAppMenu{
		TemplateID: req.TemplateID,
		Code:       req.Code,
		AppCode:    req.AppCode,
		Info:       menu.ToMenuInfoRpc(&req.MenuInfo),
	})

	return nil, err
}
