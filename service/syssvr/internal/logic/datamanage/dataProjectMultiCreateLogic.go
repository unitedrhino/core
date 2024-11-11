package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataProjectMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataProjectMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataProjectMultiCreateLogic {
	return &DataProjectMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataProjectMultiCreateLogic) DataProjectMultiCreate(in *sys.DataProjectMultiSaveReq) (*sys.Empty, error) {
	if len(in.TargetIDs) == 0 {
		return nil, errors.Parameter.WithMsg("targetIDs参数必填")
	}
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID != 0 {
		uc.ProjectID = in.ProjectID
	} else {
		in.ProjectID = uc.ProjectID
	}
	project, err := relationDB.NewProjectInfoRepo(l.ctx).FindOne(l.ctx, in.ProjectID)
	if err != nil {
		return nil, err
	}
	if !(uc.IsAdmin || uc.UserID == project.AdminUserID && in.TargetType != def.TargetRole) {
		return nil, errors.Permissions.WithMsg("只有管理员才有权限授权")
	}
	var pos []*relationDB.SysDataProject
	for _, v := range in.TargetIDs {
		pos = append(pos, &relationDB.SysDataProject{
			ProjectID:  in.ProjectID,
			TargetType: in.TargetType,
			TargetID:   v,
			AuthType:   in.AuthType,
		})
	}
	err = relationDB.NewDataProjectRepo(l.ctx).MultiInsert(l.ctx, pos)
	if err != nil && errors.Cmp(err, errors.Duplicate) {
		return nil, err
	}
	if in.TargetType == def.TargetUser {
		for _, v := range in.TargetIDs {
			cache.ClearProjectAuth(v)
		}
	}
	return &sys.Empty{}, nil
}
