package projectmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectCrudReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProjectCrudReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectCrudReadLogic {
	return &ProjectCrudReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取项目crud详情
func (l *ProjectCrudReadLogic) ProjectCrudRead(in *sys.WithID) (*sys.ProjectCrud, error) {
	po, err := relationDB.NewProjectCurlRepo(l.ctx).FindOne(l.ctx, in.Id)

	return ProjectCrudToPb(l.ctx, l.svcCtx, po), err
}
