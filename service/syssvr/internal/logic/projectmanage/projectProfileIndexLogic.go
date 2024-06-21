package projectmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectProfileIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProjectProfileIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectProfileIndexLogic {
	return &ProjectProfileIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProjectProfileIndexLogic) ProjectProfileIndex(in *sys.ProjectProfileIndexReq) (*sys.ProjectProfileIndexResp, error) {
	ret, err := relationDB.NewProjectProfileRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProjectProfileFilter{
		Codes:     in.Codes,
		ProjectID: ctxs.GetUserCtxNoNil(l.ctx).ProjectID,
	}, nil)
	return &sys.ProjectProfileIndexResp{Profiles: utils.CopySlice[sys.ProjectProfile](ret)}, err

}
