package datamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/cache"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataAreaMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataAreaMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataAreaMultiDeleteLogic {
	return &DataAreaMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataAreaMultiDeleteLogic) DataAreaMultiDelete(in *sys.DataAreaMultiDeleteReq) (*sys.Empty, error) {
	uc := ctxs.GetUserCtx(l.ctx)

	if in.ProjectID != 0 {
		uc.ProjectID = in.ProjectID
	} else {
		in.ProjectID = uc.ProjectID
	}

	if !(uc.IsAdmin) {
		pa := uc.ProjectAuth[in.ProjectID]
		if pa == nil {
			return nil, errors.Permissions.AddMsg("项目无权限")
		}
		in.TargetID = uc.UserID
		if in.TargetType != def.TargetUser {
			return nil, errors.Permissions.AddMsg("普通用户只能修改用户类型")
		}
		for _, v := range in.AreaIDs {
			if pa.Area[v] == 0 {
				return nil, errors.Permissions.AddMsg("区域无权限")
			}
		}
		//return nil, errors.Permissions.WithMsg("只有管理员才有权限授权")
	}
	err := relationDB.NewDataAreaRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DataAreaFilter{Targets: []*relationDB.Target{{Type: in.TargetType, ID: in.TargetID}}, AreaIDs: in.AreaIDs})
	if err != nil {
		return nil, err
	}
	if in.TargetType == def.TargetUser {
		cache.ClearProjectAuth(in.TargetID)
	}
	return &sys.Empty{}, err
}
