package syncJob

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取同步任务列表
func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.DeptSyncJobIndexReq) (resp *types.DeptSyncJobIndexResp, err error) {
	ret, err := l.svcCtx.DeptM.DeptSyncJobIndex(l.ctx, utils.Copy[sys.DeptSyncJobIndexReq](req))

	return utils.Copy[types.DeptSyncJobIndexResp](ret), err
}
