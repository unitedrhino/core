package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/errors"
)

func CheckModule(ctx context.Context, moduleCode string) error {
	c, err := relationDB.NewModuleInfoRepo(ctx).CountByFilter(ctx, relationDB.ModuleInfoFilter{Codes: []string{moduleCode}})
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.Parameter.AddMsgf("moduleCode not find:%v", moduleCode)
	}
	return nil
}
