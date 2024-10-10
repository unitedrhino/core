package projectmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectProfileReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProjectProfileReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectProfileReadLogic {
	return &ProjectProfileReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProjectProfileReadLogic) ProjectProfileRead(in *sys.ProjectProfileReadReq) (*sys.ProjectProfile, error) {
	ret, err := relationDB.NewProjectProfileRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ProjectProfileFilter{
		Code:      in.Code,
		ProjectID: ctxs.GetUserCtxNoNil(l.ctx).ProjectID,
	})
	if errors.Cmp(err, errors.NotFind) {
		return &sys.ProjectProfile{
			Code:   in.Code,
			Params: "",
		}, nil
	}
	return utils.Copy[sys.ProjectProfile](ret), err
}
