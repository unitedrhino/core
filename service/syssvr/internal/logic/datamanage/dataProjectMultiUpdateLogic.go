package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DataProjectMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UapDB *relationDB.DataProjectRepo
}

func NewDataProjectMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataProjectMultiUpdateLogic {
	return &DataProjectMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UapDB:  relationDB.NewDataProjectRepo(ctx),
	}
}

//func (l *DataProjectMultiUpdateLogic) DataProjectMultiUpdate(in *sys.DataProjectMultiUpdateReq) (*sys.Empty, error) {
//	uc := ctxs.GetUserCtx(l.ctx)
//	if !uc.IsAdmin {
//		return nil, errors.Permissions.WithMsg("只有管理员才有权限授权")
//	}
//	if in.TargetID == 0 {
//		return nil, errors.Parameter.AddDetail(in.TargetID).WithMsg("用户ID参数必填")
//	}
//	po, err := checkUser(l.ctx, in.TargetID)
//	if err != nil {
//		return nil, errors.Fmt(err).WithMsg("检查用户出错")
//	} else if po == nil {
//		return nil, errors.Parameter.AddDetail(err).WithMsg("检查用户不存在")
//	}
//	projects := ToAuthProjectDos(in.Projects)
//	err = l.UapDB.MultiUpdate(l.ctx, in.TargetID, projects)
//	if err != nil {
//		return nil, errors.Fmt(err).WithMsg("用户数据权限保存失败")
//	}
//
//	if in.TargetType == def.TargetUser {
//		cache.ClearProjectAuth(in.TargetID)
//	}
//	return &sys.Empty{}, nil
//}
