package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/domain/tenant"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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

	newPo := utils.Copy[relationDB.SysTenantConfig](in)
	newPo.NoDelTime = old.NoDelTime
	newPo.ID = old.ID
	if err := ctxs.IsRoot(l.ctx); err != nil {
		newPo.DeviceLimit = old.DeviceLimit
	}
	err = relationDB.NewTenantConfigRepo(l.ctx).Update(l.ctx, newPo)
	return &sys.Empty{}, err
}
