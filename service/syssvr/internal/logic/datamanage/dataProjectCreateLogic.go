package datamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataProjectCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataProjectCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataProjectCreateLogic {
	return &DataProjectCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataProjectCreateLogic) DataProjectCreate(in *sys.DataProjectSaveReq) (*sys.Empty, error) {
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
	err = relationDB.NewDataProjectRepo(l.ctx).Insert(l.ctx, &relationDB.SysDataProject{
		ProjectID:  in.ProjectID,
		TargetType: in.TargetType,
		TargetID:   in.TargetID,
		AuthType:   in.AuthType,
	})
	if err != nil && errors.Cmp(err, errors.Duplicate) {
		return nil, err
	}
	return &sys.Empty{}, nil
}
