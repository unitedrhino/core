package day

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
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

func (d *Day) Handle() error {
	var wait sync.WaitGroup
	wait.Add(1)
	utils.Go(d.ctx, func() {
		defer wait.Done()
		err := d.HandleLog()
		if err != nil {
			d.Error("%v", err)
		}
	})
	//wait.Add(1)
	//utils.Go(d.ctx, func() {
	//	defer wait.Done()
	//	err := d.HandleOss()
	//	if err != nil {
	//		d.Error("%v", err)
	//	}
	//})
	wait.Wait()
	return nil
}

//func (d *Day) HandleOss() error {
//	var now = time.Now()
//	for i := 0; i < 10; i++ {
//		now = now.AddDate(0, 0, -1)
//		path := utils.ToYYMMdd2(now.UnixMilli()) + "/"
//		err := d.svcCtx.OssClient.Delete(d.ctx, path, common.OptionKv{})
//		if err != nil {
//			d.Error("%v", err)
//		}
//	}
//	return nil
//}

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
