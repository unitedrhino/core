package projectmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectCrudDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProjectCrudDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectCrudDeleteLogic {
	return &ProjectCrudDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除项目crud
func (l *ProjectCrudDeleteLogic) ProjectCrudDelete(in *sys.WithID) (*sys.Empty, error) {
	err := relationDB.NewProjectCurlRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Empty{}, err
}
