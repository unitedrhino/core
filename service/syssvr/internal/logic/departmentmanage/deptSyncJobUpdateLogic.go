package departmentmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/domain/dept"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptSyncJobUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptSyncJobUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptSyncJobUpdateLogic {
	return &DeptSyncJobUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptSyncJobUpdateLogic) DeptSyncJobUpdate(in *sys.DeptSyncJob) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewDeptSyncJobRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.ThirdType != "" {
		old.ThirdType = in.ThirdType
	}
	if in.ThirdConfig != nil {
		newConfig := utils.Copy[relationDB.SysTenantThird](in.ThirdConfig)
		if *newConfig != *old.ThirdConfig { //这种情况需要更新下任务
			old.ThirdConfig = newConfig
		}
		if old.SyncMode == dept.SyncModeRealTime {
			err = DeptSyncAddDing(l.ctx, l.svcCtx, old)
			if err != nil {
				return nil, errors.System.WithMsgf("钉钉连接失败:%v", err.Error())
			}
		}
	}
	if in.FieldMap != nil {
		old.FieldMap = in.FieldMap
	}
	if in.SyncDeptIDs != nil {
		old.SyncDeptIDs = in.SyncDeptIDs
	}
	if in.IsAddSync != 0 {
		old.IsAddSync = in.IsAddSync
	}
	if in.SyncMode != 0 {
		old.SyncMode = in.SyncMode
	}
	err = relationDB.NewDeptSyncJobRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
