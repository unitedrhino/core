package opslogic

import (
	"context"
	"fmt"
	notifymanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/notifymanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/domain/ops"
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
	po.Status = ops.WorkOrderStatusWait
	po.RaiseUserID = ctxs.GetUserCtx(l.ctx).UserID
	err := relationDB.NewOpsFeedbackRepo(l.ctx).Insert(l.ctx, po)
	if err == nil {
		cfg, err := l.svcCtx.TenantConfigCache.GetData(l.ctx, string(po.TenantCode))
		if err != nil {
			l.Error(err)
			return &sys.WithID{Id: po.ID}, nil
		}
		if len(cfg.FeedbackNotifyUserIDs) != 0 {
			uc := ctxs.GetUserCtx(l.ctx)
			_, err = notifymanagelogic.NewMessageInfoSendLogic(l.ctx, l.svcCtx).MessageInfoSend(&sys.MessageInfoSendReq{
				UserIDs:    cfg.FeedbackNotifyUserIDs,
				IsGlobal:   0,
				NotifyCode: "feedback",
				Subject:    fmt.Sprintf("用户问题反馈: %s", po.IssueDesc),
				Body: fmt.Sprintf(`<p>反馈类型: %s</p><p>反馈内容: %s</p><p>反馈者账号: %s</p><p>联系方式: %s</p><p><a href="/app/core/#/system-manage/feedback" target="_blank">详情</a></p>`,
					po.Type, po.IssueDesc, uc.Account, po.ContactInformation),
			})
			if err != nil {
				l.Error(err)
			}
		}
	}
	return &sys.WithID{Id: po.ID}, err
}
