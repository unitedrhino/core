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

type DeptInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptInfoIndexLogic {
	return &DeptInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptInfoIndexLogic) DeptInfoIndex(in *sys.DeptInfoIndexReq) (*sys.DeptInfoIndexResp, error) {
	if in.TenantCode != "" && ctxs.IsRoot(l.ctx) != nil {
		l.ctx = ctxs.BindTenantCode(l.ctx, in.TenantCode, 0)
	}
	f := relationDB.DeptInfoFilter{
		Name:        in.Name,
		ParentID:    in.ParentID,
		Status:      in.Status,
		DingTalkIDs: in.DingTalkIDs,
	}
	repo := relationDB.NewDeptInfoRepo(l.ctx)
	total, err := repo.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := repo.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultSort())
	if err != nil {
		return nil, err
	}
	return &sys.DeptInfoIndexResp{Total: total, List: utils.CopySlice[sys.DeptInfo](pos)}, nil
}
