package projectmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoReadLogic {
	return &ProjectInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 获取项目信息详情
func (l *ProjectInfoReadLogic) ProjectInfoRead(in *sys.ProjectWithID) (*sys.ProjectInfo, error) {
	ctxs.GetUserCtx(l.ctx).AllProject = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllProject = false
	}()
	po, err := l.PiDB.FindOne(l.ctx, in.ProjectID)
	if err != nil {
		return nil, err
	}
	return transPoToPb(po), nil
}
