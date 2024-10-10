package commonlogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	f := relationDB.SlotInfoFilter{
		Code: in.Code, SubCode: in.SubCode}
	total, err := relationDB.NewSlotInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}

	list, err := relationDB.NewSlotInfoRepo(l.ctx).FindByFilter(l.ctx, f, utils.Copy[stores.PageInfo](in.Page))
	if err != nil {
		return nil, err
	}

	return &sys.SlotInfoIndexResp{List: utils.CopySlice[sys.SlotInfo](list), Total: total}, nil
}
