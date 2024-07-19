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
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID != 0 {
		uc.ProjectID = in.ProjectID
	} else {
		in.ProjectID = uc.ProjectID
	}
	if in.ProjectID != 0 {
		if uc.IsAdmin || uc.ProjectAuth[in.ProjectID] != nil {
			in.ProjectID = in.ProjectID
		}
	}

	project, err := relationDB.NewProjectInfoRepo(l.ctx).FindOne(ctxs.WithRoot(l.ctx), in.ProjectID)
	if err != nil {
		return nil, err
	}
	if !(uc.IsAdmin || (uc.UserID == project.AdminUserID && in.TargetType != def.TargetRole) ||
		(in.TargetID == uc.UserID && in.TargetType == def.TargetUser && uc.UserID != project.AdminUserID)) {
		return nil, errors.Permissions.WithMsg("只有管理员才有权限授权")
	}
	err = relationDB.NewDataProjectRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DataProjectFilter{
		ProjectID: in.GetProjectID(),
		Targets:   []*relationDB.Target{{Type: in.TargetType, ID: in.TargetID}},
	})
	return &sys.Empty{}, err
}
