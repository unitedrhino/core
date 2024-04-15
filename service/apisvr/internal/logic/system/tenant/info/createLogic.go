package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/user"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"

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

func (l *CreateLogic) Create(req *types.TenantInfoCreateReq) (*types.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	//if req.AdminUserInfo.UserName == "" {
	//	return nil, errors.Parameter.AddMsgf("需要填写管理员账号")
	//}
	resp, err := l.svcCtx.TenantRpc.TenantInfoCreate(l.ctx, &sys.TenantInfoCreateReq{Info: system.ToTenantInfoRpc(req.Info), AdminUserInfo: user.UserInfoToRpc(req.AdminUserInfo)})
	return logic.SysToWithIDTypes(resp), err
}
