package opslogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsFeedbackCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpsFeedbackCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsFeedbackCreateLogic {
	return &OpsFeedbackCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OpsFeedbackCreateLogic) OpsFeedbackCreate(in *sys.OpsFeedback) (*sys.WithID, error) {
	var po = utils.Copy[relationDB.SysOpsFeedback](in)
	po.ID = 0
	po.RaiseUserID = ctxs.GetUserCtx(l.ctx).UserID
	err := relationDB.NewOpsFeedbackRepo(l.ctx).Insert(l.ctx, po)
	return &sys.WithID{Id: po.ID}, err
}
