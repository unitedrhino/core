package opslogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsFeedbackUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpsFeedbackUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsFeedbackUpdateLogic {
	return &OpsFeedbackUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OpsFeedbackUpdateLogic) OpsFeedbackUpdate(in *sys.OpsFeedback) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllProject = true
	old, err := relationDB.NewOpsFeedbackRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Status != 0 && in.Status > old.Status {
		switch in.Status {
		case ops.WorkOrderStatusHandling:
			if old.Status == ops.WorkOrderStatusWait {
				old.Status = in.Status
			}
		case ops.WorkOrderStatusFinished:
			if old.Status == ops.WorkOrderStatusHandling {
				old.Status = in.Status
			}
		}
	}
	err = relationDB.NewOpsFeedbackRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
