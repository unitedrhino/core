package projectmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectCrudCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProjectCrudCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectCrudCreateLogic {
	return &ProjectCrudCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增项目crud
func (l *ProjectCrudCreateLogic) ProjectCrudCreate(in *sys.ProjectCrud) (*sys.WithID, error) {
	if in.Purpose == "" {
		return nil, errors.Parameter.AddMsg("purpose must not empty")
	}
	po := utils.Copy[relationDB.SysProjectCrud](in)
	po.ID = 0
	err := relationDB.NewProjectCurlRepo(l.ctx).Insert(l.ctx, po)
	return &sys.WithID{Id: po.ID}, err
}
