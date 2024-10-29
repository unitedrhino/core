package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 同步部门
func NewSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncLogic {
	return &SyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncLogic) Sync(req *types.DeptInfoSyncReq) error {
	_, err := l.svcCtx.DeptM.DeptInfoSync(l.ctx, utils.Copy[sys.DeptInfoSyncReq](req))

	return err
}
