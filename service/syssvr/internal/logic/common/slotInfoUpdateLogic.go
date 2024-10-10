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

type SlotInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSlotInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SlotInfoUpdateLogic {
	return &SlotInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SlotInfoUpdateLogic) SlotInfoUpdate(in *sys.SlotInfo) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewSlotInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	newPo := utils.Copy[relationDB.SysSlotInfo](in)
	newPo.SoftTime = old.SoftTime
	err = relationDB.NewSlotInfoRepo(l.ctx).Update(l.ctx, newPo)
	return &sys.Empty{}, err
}
