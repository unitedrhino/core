package info

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"sync"

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
