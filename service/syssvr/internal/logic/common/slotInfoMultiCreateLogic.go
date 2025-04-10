package commonlogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type SlotInfoMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSlotInfoMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SlotInfoMultiCreateLogic {
	return &SlotInfoMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SlotInfoMultiCreateLogic) SlotInfoMultiCreate(in *sys.SlotInfoMultiCreateReq) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	pos := utils.CopySlice[relationDB.SysSlotInfo](in.List)
	for _, v := range pos {
		v.ID = 0
	}
	err := relationDB.NewSlotInfoRepo(l.ctx).MultiInsert(l.ctx, pos)
	return &sys.Empty{}, err
}
