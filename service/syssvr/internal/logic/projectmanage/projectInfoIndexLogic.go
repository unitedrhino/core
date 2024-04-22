package projectmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"

	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
	var (
		list  []*sys.ProjectInfo
		total int64
		err   error
		uc    = ctxs.GetUserCtx(l.ctx)
	)
	uc.AllProject = true
	defer func() {
		uc.AllProject = false
	}()
	filter := relationDB.ProjectInfoFilter{
		ProjectIDs:  in.ProjectIDs,
		ProjectName: in.ProjectName,
	}
	if !uc.IsAdmin { //不是超管需要鉴权
		projects, err := relationDB.NewDataProjectRepo(l.ctx).FindByFilter(l.ctx,
			relationDB.DataProjectFilter{Targets: []*relationDB.Target{
				{ID: uc.UserID, Type: def.TargetUser}, {ID: uc.RoleID, Type: def.TargetRole}}}, nil)
		if err != nil {
			return nil, err
		}
		if len(projects) == 0 {
			return &sys.ProjectInfoIndexResp{}, nil
		}
		for _, v := range projects {
			filter.ProjectIDs = append(filter.ProjectIDs, v.ProjectID)
		}
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
		list = append(list, transPoToPb(po))
	}
	return &sys.ProjectInfoIndexResp{List: list, Total: total}, nil
}
