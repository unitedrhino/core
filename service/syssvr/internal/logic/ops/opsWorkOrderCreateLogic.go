package opslogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/domain/ops"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsWorkOrderCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpsWorkOrderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsWorkOrderCreateLogic {
	return &OpsWorkOrderCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 维护工单  Work Order
func (l *OpsWorkOrderCreateLogic) OpsWorkOrderCreate(in *sys.OpsWorkOrder) (*sys.WithID, error) {
	po := ToOpsWorkOrderPo(in)
	po.ID = 0
	po.Status = ops.WorkOrderStatusWait
	po.RaiseUserID = ctxs.GetUserCtx(l.ctx).UserID
	now := time.Now()
	f := relationDB.OpsWorkOrderFilter{
		StartTime: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
		EndTime:   time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local),
	}
	todayCount, err := relationDB.NewOpsWorkOrderRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	po.Number = fmt.Sprintf("DMWO%04d%02d%02d%04d", now.Year(), now.Month(), now.Day(), todayCount+1)
	err = relationDB.NewOpsWorkOrderRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}
