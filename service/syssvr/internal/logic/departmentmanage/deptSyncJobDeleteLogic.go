package departmentmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/domain/dept"
	"gitee.com/unitedrhino/core/service/syssvr/internal/event/deptSync"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptSyncJobDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptSyncJobDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptSyncJobDeleteLogic {
	return &DeptSyncJobDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptSyncJobDeleteLogic) DeptSyncJobDelete(in *sys.WithID) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewDeptSyncJobRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if old.SyncMode == dept.SyncModeRealTime {
		err = deptSync.NewDeptSync(l.ctx, l.svcCtx).DelDing(old)
		if err != nil {
			return nil, err
		}
	}
	err = relationDB.NewDeptSyncJobRepo(l.ctx).Delete(l.ctx, in.Id)

	return &sys.Empty{}, err
}
