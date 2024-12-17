package deptSync

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/share/clients/dingClient"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
)

type DeptSync struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	log    logx.Logger
}

func NewDeptSync(ctx context.Context, svcCtx *svc.ServiceContext) *DeptSync {
	return &DeptSync{
		ctx:    ctx,
		log:    logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (d DeptSync) HandleDing(df *payload.DataFrame) error {
	fmt.Println(df)
	return nil
}

func (d DeptSync) AddDing(po *relationDB.SysDeptSyncJob) error {
	d.svcCtx.DingStreamMapMutex.Lock()
	defer d.svcCtx.DingStreamMapMutex.Unlock()
	cli1 := d.svcCtx.DingStreamMap[string(po.TenantCode)]
	if cli1 != nil {
		defer cli1.Close()
	}
	cli := dingClient.NewDingStream(po.ThirdConfig.AppKey, po.ThirdConfig.AppSecret)
	cli.RegisterAllEventRouter(func(c context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
		ctx, span := ctxs.StartSpan(c, "dingStreamEvent", "")
		defer span.End()
		startTime := timex.Now()
		duration := timex.Since(startTime)
		ctx = ctxs.WithRoot(ctx)
		ctx = ctxs.BindTenantCode(ctx, string(po.TenantCode), 0)
		err := NewDeptSync(ctx, d.svcCtx).HandleDing(df)
		logx.WithContext(ctx).WithDuration(duration).Infof(
			"subscribeDingStream df:%v err:%v",
			utils.Fmt(df), err)
		return event.NewSuccessResponse()
	})
	err := cli.Start(context.Background())
	if err != nil {
		logx.Error(err)
		return err
	}
	d.svcCtx.DingStreamMap[string(po.TenantCode)] = cli
	return nil
}

func (d DeptSync) DelDing(po *relationDB.SysDeptSyncJob) error {
	d.svcCtx.DingStreamMapMutex.Lock()
	defer d.svcCtx.DingStreamMapMutex.Unlock()
	cli1 := d.svcCtx.DingStreamMap[string(po.TenantCode)]
	if cli1 != nil {
		cli1.Close()
		delete(d.svcCtx.DingStreamMap, string(po.TenantCode))
	}
	return nil
}
