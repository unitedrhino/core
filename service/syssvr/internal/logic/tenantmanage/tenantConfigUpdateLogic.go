package tenantmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/domain/tenant"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantConfigUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantConfigUpdateLogic {
	return &TenantConfigUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantConfigUpdateLogic) TenantConfigUpdate(in *sys.TenantConfig) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	if ctxs.IsRoot(l.ctx) == nil {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	old, err := relationDB.NewTenantConfigRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantConfigFilter{TenantCode: in.TenantCode})
	if err != nil {
		return nil, err
	}
	if in.RegisterAutoCreateProject != nil {
		oldPMap, maxProjectID := tenant.RegisterAutoCreateProjectToMap(old.RegisterAutoCreateProject)
		for _, v := range in.RegisterAutoCreateProject {
			if v.Id == 0 {
				maxProjectID++
				v.Id = maxProjectID
			}
			oldP := oldPMap[v.Id]
			if oldP == nil {
				oldP = &tenant.RegisterAutoCreateProject{
					ID:           v.Id,
					ProjectName:  v.ProjectName,
					IsSysCreated: v.IsSysCreated,
				}
			}
			maxAreaID := oldP.MaxAreaID
			for _, a := range v.Areas {
				if a.Id == 0 {
					maxAreaID++
					a.Id = maxAreaID
				}
				oldA := oldP.AreaMap[a.Id]
				if a.IsUpdateAreaImg == true && a.AreaImg != "" {
					nwePath := oss.GenCommonFilePath(l.svcCtx.Config.Name, oss.BusinessArea, oss.SceneHeadIng, oss.GetFileNameWithPath(a.AreaImg))
					path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(a.AreaImg, nwePath)
					if err != nil {
						l.Error(err)
					} else {
						a.AreaImg = path
					}
				}
				if !a.IsUpdateAreaImg && oldA != nil { //更新类型
					a.AreaImg = oldA.AreaImg
				}
			}
		}
		old.RegisterAutoCreateProject = utils.CopySlice[tenant.RegisterAutoCreateProject](in.RegisterAutoCreateProject)
	}
	if ctxs.IsRoot(l.ctx) == nil {
		if in.DeviceLimit != nil {
			old.DeviceLimit = in.DeviceLimit.GetValue()
		}
		if in.LoginLogKeepDays != nil {
			old.LoginLogKeepDays = in.LoginLogKeepDays.GetValue()
		}
		if in.OperLogKeepDays != nil {
			old.OperLogKeepDays = in.OperLogKeepDays.GetValue()
		}
		if in.UserLimit != nil {
			old.UserLimit = in.UserLimit.GetValue()
		}
		if in.ProjectLimit != nil {
			old.ProjectLimit = in.ProjectLimit.GetValue()
		}
	}
	if in.FeedbackNotifyUserIDs != nil {
		old.FeedbackNotifyUserIDs = in.FeedbackNotifyUserIDs
	}
	if in.RegisterRoleID != 0 {
		old.RegisterRoleID = in.RegisterRoleID
	}
	if in.CheckUserDelete != 0 {
		old.CheckUserDelete = in.CheckUserDelete
	}

	err = relationDB.NewTenantConfigRepo(l.ctx).Update(l.ctx, old)
	if err == nil {
		l.svcCtx.TenantConfigCache.SetData(l.ctx, string(old.TenantCode), nil)
	}
	return &sys.Empty{}, err
}
