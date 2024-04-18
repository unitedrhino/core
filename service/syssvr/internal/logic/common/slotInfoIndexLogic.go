package commonlogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type SlotInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSlotInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SlotInfoIndexLogic {
	return &SlotInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SlotInfoIndexLogic) SlotInfoIndex(in *sys.SlotInfoIndexReq) (*sys.SlotInfoIndexResp, error) {
	ret := l.svcCtx.Slot.Get(l.ctx, in.Code, in.SubCode)
	if ret == nil {
		return &sys.SlotInfoIndexResp{}, errors.NotFind
	}
	return &sys.SlotInfoIndexResp{Slots: utils.CopySlice[sys.SlotInfo](ret)}, nil
}
