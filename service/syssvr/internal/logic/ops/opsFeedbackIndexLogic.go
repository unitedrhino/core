package opslogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsFeedbackIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpsFeedbackIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsFeedbackIndexLogic {
	return &OpsFeedbackIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OpsFeedbackIndexLogic) OpsFeedbackIndex(in *sys.OpsFeedbackIndexReq) (*sys.OpsFeedbackIndexResp, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	if err := ctxs.IsRoot(l.ctx); err == nil && in.IsAllTenant {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	ctxs.GetUserCtx(l.ctx).AllProject = true
	f := relationDB.OpsFeedbackFilter{
		TenantCode: in.TenantCode,
		ProjectID:  in.ProjectID,
		Type:       in.Type,
		Status:     in.Status,
	}
	total, err := relationDB.NewOpsFeedbackRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	list, err := relationDB.NewOpsFeedbackRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "createdTime",
		Sort:  stores.OrderDesc,
	}))
	if err != nil {
		return nil, err
	}
	return &sys.OpsFeedbackIndexResp{List: utils.CopySlice[sys.OpsFeedback](list), Total: total}, nil
}
