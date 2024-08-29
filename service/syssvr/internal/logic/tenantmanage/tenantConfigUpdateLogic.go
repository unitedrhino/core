package tenantmanagelogic

import (
	"context"
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
	for _, v := range in.RegisterAutoCreateProject {
		for _, a := range v.Areas {
			if a.IsUpdateAreaImg == true && a.AreaImg != "" {
				nwePath := oss.GenCommonFilePath(l.svcCtx.Config.Name, oss.BusinessArea, oss.SceneHeadIng, oss.GetFileNameWithPath(a.AreaImg))
				path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(a.AreaImg, nwePath)
				if err != nil {
					l.Error(err)
				} else {
					a.AreaImg = path
				}
			}
		}
	}
	newPo := utils.Copy[relationDB.SysTenantConfig](in)
	newPo.NoDelTime = old.NoDelTime
	newPo.ID = old.ID
	err = relationDB.NewTenantConfigRepo(l.ctx).Update(l.ctx, newPo)
	return &sys.Empty{}, err
}
