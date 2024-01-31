package self

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaApplyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAreaApplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaApplyCreateLogic {
	return &AreaApplyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaApplyCreateLogic) AreaApplyCreate(req *types.UserAreaApplyCreateReq) error {
	_, err := l.svcCtx.UserRpc.UserAreaApplyCreate(l.ctx, &sys.UserAreaApplyCreateReq{
		AreaID:   req.AreaID,
		AuthType: req.AuthType,
	})
	return err
}
