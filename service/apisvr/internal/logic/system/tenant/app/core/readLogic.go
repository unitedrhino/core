package core

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.WithCode) (resp *types.AppCore, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppRead(l.ctx, &sys.TenantAppWithIDOrCode{AppCode: req.Code})
	if err != nil {
		return nil, err
	}
	ret2, err := l.svcCtx.TenantRpc.TenantInfoRead(l.ctx, &sys.WithIDCode{})
	if err != nil {
		return nil, err
	}
	r := utils.Copy[types.AppCore](ret)
	r.Tenant = utils.Copy[types.TenantCore](ret2)
	return r, nil
}
