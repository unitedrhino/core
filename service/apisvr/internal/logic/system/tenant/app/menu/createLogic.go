package menu

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/module/menu"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.TenantAppMenu) (resp *types.WithID, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppMenuCreate(l.ctx, &sys.TenantAppMenu{
		TemplateID: req.TemplateID,
		Code:       req.Code,
		AppCode:    req.AppCode,
		Info:       menu.ToMenuInfoRpc(&req.MenuInfo),
	})
	if err != nil {
		return nil, err
	}
	return &types.WithID{
		ID: ret.Id,
	}, nil
}
