package departmentmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptSyncJobIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptSyncJobIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptSyncJobIndexLogic {
	return &DeptSyncJobIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptSyncJobIndexLogic) DeptSyncJobIndex(in *sys.DeptSyncJobIndexReq) (*sys.DeptSyncJobIndexResp, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	f := relationDB.DeptSyncJobFilter{
		Direction: in.Direction,
		ThirdType: in.ThirdType,
		SyncMode:  in.SyncMode,
	}
	repo := relationDB.NewDeptSyncJobRepo(l.ctx)
	total, err := repo.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := repo.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	return &sys.DeptSyncJobIndexResp{Total: total, List: utils.CopySlice[sys.DeptSyncJob](pos)}, nil
}
