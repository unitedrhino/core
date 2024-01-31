package menu

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MulitUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMulitUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MulitUpdateLogic {
	return &MulitUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MulitUpdateLogic) MulitUpdate(req *types.RoleMenuMultiUpdateReq) error {
	resp, err := l.svcCtx.RoleRpc.RoleMenuMultiUpdate(l.ctx, &sys.RoleMenuMultiUpdateReq{
		Id:         req.ID,
		AppCode:    req.AppCode,
		ModuleCode: req.ModuleCode,
		MenuIDs:    req.MenuIDs,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.RoleMenuUpdate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s.rpc.RoleMenuUpdate return nil req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("RoleMenuUpdate rpc return nil")
	}
	return nil
}
