package projectmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoIndexLogic {
	return &ProjectInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 获取项目信息列表
func (l *ProjectInfoIndexLogic) ProjectInfoIndex(in *sys.ProjectInfoIndexReq) (*sys.ProjectInfoIndexResp, error) {
	if in.TenantCode != "" && ctxs.IsRoot(l.ctx) == nil {
		l.ctx = ctxs.BindTenantCode(l.ctx, in.TenantCode, 0)
	}
	var (
		list  []*sys.ProjectInfo
		total int64
		err   error
		uc    = ctxs.GetUserCtx(l.ctx)
	)
	uc.ProjectID = 0 //获取全部项目
	filter := relationDB.ProjectInfoFilter{
		ProjectIDs:   in.ProjectIDs,
		ProjectName:  in.ProjectName,
		WithTopAreas: in.WithTopAreas,
	}
	if in.IsGetAll && uc.IsAdmin {
		l.ctx = ctxs.WithAllProject(l.ctx)
	}
	total, err = l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	poArr, err := l.PiDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}

	list = make([]*sys.ProjectInfo, 0, len(poArr))
	for _, po := range poArr {
		list = append(list, ProjectInfoToPb(l.ctx, l.svcCtx, po))
	}
	return &sys.ProjectInfoIndexResp{List: list, Total: total}, nil
}
