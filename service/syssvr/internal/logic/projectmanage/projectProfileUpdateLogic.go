package projectmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
		old = &relationDB.SysProjectProfile{ProjectID: stores.ProjectID(projectID), Code: in.Code}
	}
	old.Params = in.Params
	err = relationDB.NewProjectProfileRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}