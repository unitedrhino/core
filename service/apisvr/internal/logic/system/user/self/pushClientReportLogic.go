package self

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushClientReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPushClientReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushClientReportLogic {
	return &PushClientReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushClientReportLogic) PushClientReport(req *types.UserPushClientReportReq) error {
	_, err := l.svcCtx.UserRpc.UserPushClientReport(l.ctx, utils.Copy[sys.UserPushClientReportReq](req))
	return err
}
