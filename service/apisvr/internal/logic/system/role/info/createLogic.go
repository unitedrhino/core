package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

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

func (l *CreateLogic) Create(req *types.RoleInfo) (*types.WithID, error) {
	if req.Status == 0 {
		req.Status = 1
	}
	resp, err := l.svcCtx.RoleRpc.RoleInfoCreate(l.ctx, utils.Copy[sys.RoleInfo](req))
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.RoleCreate req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	return logic.SysToWithIDTypes(resp), nil
}
