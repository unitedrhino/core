package datamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataProjectDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataProjectDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataProjectDeleteLogic {
	return &DataProjectDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataProjectDeleteLogic) DataProjectDelete(in *sys.DataProjectDeleteReq) (*sys.Empty, error) {
	if in.TargetID == 0 {
		return nil, errors.Parameter.AddDetail(in.TargetID).WithMsg("TargetID参数必填")
	}
	if in.ProjectID == 0 {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("项目id参数必填")
	}
	project, err := relationDB.NewDataProjectRepo(l.ctx).FindOne(l.ctx, in.TargetType, in.TargetID, in.ProjectID)
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
	}
	uc := ctxs.GetUserCtx(l.ctx)
	if !(uc.IsAdmin || uc.UserID == project.ProjectID && in.TargetType != def.TargetRole) {
		return nil, errors.Permissions.WithMsg("只有管理员才有权限授权")
	}
	err = relationDB.NewDataProjectRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DataProjectFilter{
		ProjectID: in.GetProjectID(),
		Targets:   []*relationDB.Target{{Type: in.TargetType, ID: in.TargetID}},
	})
	return &sys.Empty{}, err
}
