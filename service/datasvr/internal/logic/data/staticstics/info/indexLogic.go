package info

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"github.com/maypok86/otter"
	"sync"
	"time"

	"gitee.com/unitedrhino/core/service/datasvr/internal/svc"
	"gitee.com/unitedrhino/core/service/datasvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var cache otter.Cache[string, *types.StaticsticsInfoReadResp]

func init() {
	c, err := otter.MustBuilder[string, *types.StaticsticsInfoReadResp](10_000).
		CollectStats().
		Cost(func(key string, value *types.StaticsticsInfoReadResp) uint32 {
			return 1
		}).
		WithTTL(5 * time.Minute).
		Build()
	logx.Must(err)
	cache = c
}

func (l *IndexLogic) Index(req *types.StaticsticsInfoIndexReq) (resp *types.StaticsticsInfoIndexResp, err error) {
	var wait sync.WaitGroup
	var rets = make([][]map[string]interface{}, len(req.Finds))
	read := NewReadLogic(l.ctx, l.svcCtx)
	for i, v := range req.Finds {
		i1, v1 := i, v
		wait.Add(1)
		utils.Go(l.ctx, func() {
			defer wait.Done()
			ret, er := read.Handle(v1)
			if er != nil {
				l.Error(er)
				return
			}
			rets[i1] = ret.List
		})
	}
	wait.Wait()
	return &types.StaticsticsInfoIndexResp{Lists: rets}, nil
}
