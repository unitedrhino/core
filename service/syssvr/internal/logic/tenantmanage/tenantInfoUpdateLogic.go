package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoUpdateLogic {
	return &TenantInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新区域
func (l *TenantInfoUpdateLogic) TenantInfoUpdate(in *sys.TenantInfo) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	repo := relationDB.NewTenantInfoRepo(l.ctx)
	old, err := repo.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.AdminUserID != 0 && in.AdminUserID != old.AdminUserID { //只有default的超管才有权限修改管理员
		err := logic.IsSupperAdmin(l.ctx, def.TenantCodeDefault)
		if err != nil {
			return nil, err
		}
		old.AdminUserID = in.AdminUserID
	}
	if in.BaseUrl != "" {
		old.BaseUrl = in.BaseUrl
	}
	if in.LogoUrl != "" {
		old.LogoUrl = in.LogoUrl
	}
	if in.Desc != nil {
		old.Desc = utils.ToEmptyString(in.Desc)
	}
	err = repo.Update(l.ctx, old)
	err = caches.SetTenant(l.ctx, logic.ToTenantInfoCache(old))
	if err != nil {
		l.Error(err)
	}
	err = l.svcCtx.TenantCache.SetData(l.ctx, old.Code, logic.ToTenantInfoCache(old))
	if err != nil {
		l.Error(err)
	}
	return &sys.Response{}, err
}
