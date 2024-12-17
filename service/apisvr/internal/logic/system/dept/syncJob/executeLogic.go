package syncJob

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 执行同步任务
func NewExecuteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteLogic {
	return &ExecuteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExecuteLogic) Execute(req *types.DeptSyncJobExecuteReq) error {
	_, err := l.svcCtx.DeptM.DeptSyncJobExecute(l.ctx, utils.Copy[sys.DeptSyncJobExecuteReq](req))

	return err
}
