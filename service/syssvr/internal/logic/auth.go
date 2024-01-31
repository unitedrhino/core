package logic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
)

// 租户号为空则为该用户的租户
func IsSupperAdmin(ctx context.Context, tenantCode string) error {
	uc := ctxs.GetUserCtx(ctx)
	if tenantCode == "" {
		tenantCode = uc.TenantCode
	}
	if uc.TenantCode != tenantCode {
		return errors.Permissions.AddMsgf("只有%s的超管才有权限", tenantCode)
	}
	ti, err := relationDB.NewTenantInfoRepo(ctx).FindOneByFilter(ctx, relationDB.TenantInfoFilter{Codes: []string{uc.TenantCode}})
	if err != nil {
		return err
	}
	if ti.AdminUserID != uc.UserID {
		return errors.Permissions.AddMsgf("只有%s的超管才有权限", tenantCode)
	}
	return nil
}
