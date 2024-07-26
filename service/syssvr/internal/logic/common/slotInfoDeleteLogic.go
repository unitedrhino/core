package commonlogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type SlotInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSlotInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SlotInfoDeleteLogic {
	return &SlotInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SlotInfoDeleteLogic) SlotInfoDelete(in *sys.WithID) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	err := relationDB.NewSlotInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Empty{}, err
}
