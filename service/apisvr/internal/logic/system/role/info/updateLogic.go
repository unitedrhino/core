package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
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

func (l *UpdateLogic) Update(req *types.RoleInfo) error {
	resp, err := l.svcCtx.RoleRpc.RoleInfoUpdate(l.ctx, utils.Copy[sys.RoleInfo](req))
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.RoleUpdate req=%v err=%v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s.rpc.RoleUpdate return nil req=%v", utils.FuncName(), req)
		return errors.System.AddDetail("RoleUpdate rpc return nil")
	}
	return nil
}
