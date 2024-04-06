package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoDeleteLogic {
	return &TenantInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除区域
func (l *TenantInfoDeleteLogic) TenantInfoDelete(in *sys.WithIDCode) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	f := relationDB.TenantInfoFilter{ID: in.Id}
	if in.Code != "" {
		f.Codes = []string{in.Code}
	}
	conn := stores.GetTenantConn(l.ctx)
	var (
		ti  *relationDB.SysTenantInfo
		err error
	)

	err = conn.Transaction(func(tx *gorm.DB) error {
		tir := relationDB.NewTenantInfoRepo(tx)
		ti, err = tir.FindOneByFilter(l.ctx, f)
		if err != nil {
			return err
		}
		if ti.Code == def.TenantCodeDefault {
			return errors.Parameter.AddMsg("默认租户不允许删除")
		}
		err = relationDB.NewTenantAppRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppModuleRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppModuleFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAccessRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAccessFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppMenuRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppMenuFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewUserInfoRepo(tx).DeleteByFilter(l.ctx, relationDB.UserInfoFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewRoleModuleRepo(tx).DeleteByFilter(l.ctx, relationDB.RoleModuleFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewRoleAccessRepo(tx).DeleteByFilter(l.ctx, relationDB.RoleAccessFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewRoleAppRepo(tx).DeleteByFilter(l.ctx, relationDB.RoleAppFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewRoleMenuRepo(tx).DeleteByFilter(l.ctx, relationDB.RoleMenuFilter{TenantCode: ti.Code})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantInfoRepo(l.ctx).DeleteByFilter(l.ctx, f)
		if err != nil {
			return err
		}
		err = caches.DelTenant(l.ctx, ti.Code)
		if err != nil {
			l.Error(err)
		}
		err = l.svcCtx.TenantCache.SetData(l.ctx, ti.Code, nil)
		if err != nil {
			l.Error(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &sys.Empty{}, err
}
