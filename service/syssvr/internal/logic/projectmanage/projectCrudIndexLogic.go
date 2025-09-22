package projectmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectCrudIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProjectCrudIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectCrudIndexLogic {
	return &ProjectCrudIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取项目crud列表
func (l *ProjectCrudIndexLogic) ProjectCrudIndex(in *sys.ProjectCrudIndexReq) (*sys.ProjectCrudIndexResp, error) {
	if in.Purpose == "" {
		return &sys.ProjectCrudIndexResp{}, errors.Parameter.AddMsg("purpose must not empty")
	}
	f := relationDB.ProjectCurlFilter{
		Purpose: in.Purpose,
		Params:  utils.CopyMap[stores.Compare](in.Params),
	}
	total, err := relationDB.NewProjectCurlRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := relationDB.NewProjectCurlRepo(l.ctx).FindByFilter(l.ctx, f, utils.Copy[stores.PageInfo](in.Page).WithDefaultSort())
	if err != nil {
		return nil, err
	}
	return &sys.ProjectCrudIndexResp{List: ProjectCrudsToPb(l.ctx, l.svcCtx, pos), Total: total}, nil
}
