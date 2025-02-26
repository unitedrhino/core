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

type DataAreaMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UaaDB *relationDB.DataAreaRepo
	UapDB *relationDB.DataProjectRepo
}

func NewDataAreaMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataAreaMultiUpdateLogic {
	return &DataAreaMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UaaDB:  relationDB.NewDataAreaRepo(ctx),
		UapDB:  relationDB.NewDataProjectRepo(ctx),
	}
}

func (l *DataAreaMultiUpdateLogic) DataAreaMultiUpdate(in *sys.DataAreaMultiUpdateReq) (*sys.Empty, error) {
	if in.TargetID == 0 {
		return nil, errors.Parameter.AddDetail(in.TargetID).WithMsg("TargetID参数必填")
	}
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID != 0 {
		uc.ProjectID = in.ProjectID
	} else {
		in.ProjectID = uc.ProjectID
	}
	project, err := l.UapDB.FindOne(l.ctx, in.TargetType, in.TargetID, in.ProjectID)
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
	}
	if !(uc.IsAdmin || uc.ProjectID == project.ProjectID) {
		return nil, errors.Permissions.WithMsg("只有管理员才有权限授权")
	}
	//po, err := checkUser(l.ctx, in.TargetID)
	//if err != nil {
	//	return nil, errors.Fmt(err).WithMsg("检查用户出错")
	//} else if po == nil {
	//	return nil, errors.Parameter.AddDetail(err).WithMsg("检查用户不存在")
	//}
	areas := ToAuthAreaDos(l.ctx, l.svcCtx, in.Areas)
	err = l.UaaDB.MultiUpdate(l.ctx, &relationDB.Target{Type: in.TargetType, ID: in.TargetID}, in.ProjectID, areas)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("用户数据权限保存失败")
	}
	if len(areas) == 0 && project != nil { //如果把项目下所有区域权限取消了,则项目权限默认也取消
		err = l.UapDB.Delete(l.ctx, in.TargetType, in.TargetID, in.ProjectID)
		if err != nil {
			l.Error(err)
			return nil, err
		}
		//InitCacheUserAuthProject(l.ctx, in.TargetID)
	}
	if len(areas) != 0 && project == nil {
		err = l.UapDB.Insert(l.ctx, &relationDB.SysDataProject{TargetType: def.TargetUser, TargetID: in.TargetID, ProjectID: in.ProjectID, AuthType: def.AuthRead})
		if err != nil {
			l.Error(err)
			return nil, err
		}
		//InitCacheUserAuthProject(l.ctx, in.TargetID)
	}
	if in.TargetType == def.TargetUser {
		cache.ClearProjectAuth(in.TargetID)
		ProjectUserCount(l.ctx, in.ProjectID)
	}
	return &sys.Empty{}, nil
}
