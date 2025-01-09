package day

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type Day struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewDaySync(ctx context.Context, svcCtx *svc.ServiceContext) *Day {
	return &Day{
		ctx:    ctx,
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (d *Day) HandleLog() error {
	cs, err := relationDB.NewTenantConfigRepo(d.ctx).FindByFilter(d.ctx, relationDB.TenantConfigFilter{}, nil)
	if err != nil {
		return err
	}
	for _, c := range cs {
		if c.OperLogKeepDays != 0 {
			err := relationDB.NewOperLogRepo(d.ctx).DeleteByFilter(d.ctx, relationDB.OperLogFilter{
				TenantCode: string(c.TenantCode), CreateTime: stores.CmpLt(time.Now().Add(-time.Hour * 24 * time.Duration(c.OperLogKeepDays)))})
			if err != nil {
				d.Error(c, err)
			}
		}
		if c.LoginLogKeepDays != 0 {
			err := relationDB.NewLoginLogRepo(d.ctx).DeleteByFilter(d.ctx, relationDB.LoginLogFilter{
				TenantCode: string(c.TenantCode), CreateTime: stores.CmpLt(time.Now().Add(-time.Hour * 24 * time.Duration(c.OperLogKeepDays)))})
			if err != nil {
				d.Error(c, err)
			}
		}
	}
	return nil
}
