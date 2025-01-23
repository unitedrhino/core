package projectmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectProfileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProjectProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectProfileUpdateLogic {
	return &ProjectProfileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProjectProfileUpdateLogic) ProjectProfileUpdate(in *sys.ProjectProfile) (*sys.Empty, error) {
	projectID := ctxs.GetUserCtxNoNil(l.ctx).ProjectID
	old, err := relationDB.NewProjectProfileRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ProjectProfileFilter{
		Code:      in.Code,
		ProjectID: projectID,
	})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		old = &relationDB.SysProjectProfile{ProjectID: dataType.ProjectID(projectID), Code: in.Code}
	}
	old.Params = in.Params
	err = relationDB.NewProjectProfileRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
