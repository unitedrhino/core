package departmentmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/domain/dept"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptSyncJobCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptSyncJobCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptSyncJobCreateLogic {
	return &DeptSyncJobCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptSyncJobCreateLogic) DeptSyncJobCreate(in *sys.DeptSyncJob) (*sys.WithID, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	in.CreatedTime = 0
	in.Id = 0
	if in.Direction != dept.SyncDirectionTo { //从上游同步过来的需要检查,一个租户只能有一个,避免错乱
		total, err := relationDB.NewDeptSyncJobRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeptSyncJobFilter{Direction: dept.SyncDirectionFrom})
		if err != nil {
			return nil, err
		}
		if total > 0 {
			return nil, errors.Parameter.AddMsg("同时只能存在一个上游同步任务")
		}
	}

	po := utils.Copy[relationDB.SysDeptSyncJob](in)
	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewDeptSyncJobRepo(l.ctx).Insert(l.ctx, po)
		if err != nil {
			return err
		}
		if in.SyncMode == dept.SyncModeRealTime {
			err = DeptSyncAddDing(l.ctx, l.svcCtx, po)
		}
		return err
	})

	return &sys.WithID{Id: po.ID}, err
}
