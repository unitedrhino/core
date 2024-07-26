package commonlogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type SlotInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSlotInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SlotInfoReadLogic {
	return &SlotInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SlotInfoReadLogic) SlotInfoRead(in *sys.WithID) (*sys.SlotInfo, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po, err := relationDB.NewSlotInfoRepo(l.ctx).FindOne(l.ctx, in.Id)

	return utils.Copy[sys.SlotInfo](po), err
}
