package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataProjectMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataProjectMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataProjectMultiDeleteLogic {
	return &DataProjectMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataProjectMultiDeleteLogic) DataProjectMultiDelete(in *sys.DataProjectMultiDeleteReq) (*sys.Empty, error) {
	if len(in.TargetIDs) == 0 {
		return nil, errors.Parameter.WithMsg("TargetID参数必填")
	}
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID != 0 {
		if uc.IsAdmin || uc.ProjectAuth[in.ProjectID] != nil {
			uc.ProjectID = in.ProjectID
		} else {
			return nil, errors.Permissions
		}
	} else {
		in.ProjectID = uc.ProjectID
	}

	project, err := relationDB.NewProjectInfoRepo(l.ctx).FindOne(ctxs.WithRoot(l.ctx), in.ProjectID)
	if err != nil {
		return nil, err
	}
	if !(uc.IsAdmin || project.AdminUserID == uc.UserID) {
		return nil, errors.Permissions.WithMsg("只有管理员才有权限授权")
	}
	err = relationDB.NewDataProjectRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DataProjectFilter{
		ProjectID:  in.GetProjectID(),
		TargetIDs:  in.TargetIDs,
		TargetType: in.TargetType,
	})
	if in.TargetType == def.TargetUser {
		for _, v := range in.TargetIDs {
			cache.ClearProjectAuth(v)
		}
		ProjectUserCount(l.ctx, in.ProjectID)
	}
	return &sys.Empty{}, nil
}
