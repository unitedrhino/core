package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/stores"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppDeleteLogic {
	return &TenantAppDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppDeleteLogic) TenantAppDelete(in *sys.TenantAppWithIDOrCode) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	f := relationDB.TenantAppFilter{
		TenantCode: in.Code,
		AppCodes:   []string{in.AppCode},
	}
	if in.AppCode != "" {
		f.AppCodes = []string{in.AppCode}
	}
	if in.Id != 0 {
		f.IDs = []int64{in.Id}
	}

	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewTenantAppRepo(tx).DeleteByFilter(l.ctx, f)
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppModuleRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppModuleFilter{
			TenantCode: in.Code,
			AppCode:    in.AppCode,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppMenuRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppMenuFilter{
			TenantCode: in.Code,
			AppCode:    in.AppCode,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return &sys.Response{}, err
}
