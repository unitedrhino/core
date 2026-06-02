package config

import (
	"context"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TestSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestSendLogic {
	return &TestSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestSendLogic) TestSend(req *types.NotifyConfigTestSendReq) error {
	_, err := l.svcCtx.NotifyM.NotifyConfigSend(l.ctx, &sys.NotifyConfigSendReq{
		UserIDs:    req.UserIDs,
		Accounts:   req.Accounts,
		NotifyCode: req.NotifyCode,
		Type:       req.Type,
		Params:     req.Params,
		TemplateID: req.TemplateID,
	})
	return err
}
