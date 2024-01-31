package module

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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

func (l *MulitUpdateLogic) MulitUpdate(req *types.RoleModuleMultiUpdateReq) error {
	_, err := l.svcCtx.RoleRpc.RoleModuleMultiUpdate(l.ctx, &sys.RoleModuleMultiUpdateReq{
		Id: req.ID, AppCode: req.AppCode, ModuleCodes: req.ModuleCodes})
	return err
}
