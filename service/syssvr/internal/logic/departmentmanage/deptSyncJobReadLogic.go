package departmentmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptSyncJobReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptSyncJobReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptSyncJobReadLogic {
	return &DeptSyncJobReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptSyncJobReadLogic) DeptSyncJobRead(in *sys.DeptSyncJobReadReq) (*sys.DeptSyncJob, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	po, err := relationDB.NewDeptSyncJobRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[sys.DeptSyncJob](po), err
}
