package projectmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/errors"
)

func checkProject(ctx context.Context, productID int64) (*relationDB.SysProjectInfo, error) {
	po, err := relationDB.NewProjectInfoRepo(ctx).FindOne(ctx, productID)
	if err == nil {
		return po, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}
