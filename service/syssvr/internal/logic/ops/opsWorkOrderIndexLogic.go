package opslogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type OpsWorkOrderIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpsWorkOrderIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsWorkOrderIndexLogic {
	return &OpsWorkOrderIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OpsWorkOrderIndexLogic) OpsWorkOrderIndex(in *sys.OpsWorkOrderIndexReq) (*sys.OpsWorkOrderIndexResp, error) {
	f := relationDB.OpsWorkOrderFilter{Status: in.Status, AreaID: in.AreaID, Type: in.Type, Number: in.Number}
	total, err := relationDB.NewOpsWorkOrderRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	list, err := relationDB.NewOpsWorkOrderRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "createdTime",
		Sort:  stores.OrderDesc,
	}))
	if err != nil {
		return nil, err
	}
	return &sys.OpsWorkOrderIndexResp{List: utils.CopySlice[sys.OpsWorkOrder](list), Total: total}, nil
}
