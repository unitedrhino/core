package self

import (
	"context"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 租户管理员用户注册(配置开启了才能用)
func NewTaRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaRegisterLogic {
	return &TaRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaRegisterLogic) TaRegister(req *types.UserTaRegisterReq) error {
	_, err := l.svcCtx.UserRpc.UserTaRegister(l.ctx, utils.Copy[sys.UserTaRegisterReq](req))
	return err
}
