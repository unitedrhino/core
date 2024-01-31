package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/errors"
)

func CheckModule(ctx context.Context, tenantCode, appCode, moduleCode string) error {
	c, err := relationDB.NewTenantAppModuleRepo(ctx).CountByFilter(ctx, relationDB.TenantAppModuleFilter{TenantCode: tenantCode, AppCodes: []string{appCode}, ModuleCodes: []string{moduleCode}})
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.Parameter.AddMsgf("moduleCode not find:%v", moduleCode)
	}
	return nil
}
