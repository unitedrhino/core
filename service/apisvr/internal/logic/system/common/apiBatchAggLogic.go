package common

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
	"time"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiBatchAggLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量聚合接口请求
func NewApiBatchAggLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiBatchAggLogic {
	return &ApiBatchAggLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiBatchAggLogic) ApiBatchAgg(r *http.Request, req *types.ApiBatchAggReq) (resp *types.ApiBatchAggResp, err error) {
	resp = &types.ApiBatchAggResp{}
	var egg errgroup.Group
	var m sync.Mutex
	for _, v := range req.Reqs {
		egg.Go(func() error {
			defer utils.Recover(l.ctx)
			var rets []map[string]interface{}
			gr := gorequest.New().Retry(2, time.Second*5)
			gr.Post(fmt.Sprintf("http://localhost:%d%s", l.svcCtx.Config.Port, v.Uri))
			for k, v := range r.Header {
				gr.Set(k, v[0])
			}
			for _, body := range v.Bodys {
				rsp, b, errs := gr.Type("json").Send(body).End()
				if errs != nil {
					err = errors.System.AddDetail(errs)
				}
				if rsp.StatusCode != 200 {
					return errors.System.AddDetail(b)
				}
				var ret = map[string]interface{}{}
				err = json.Unmarshal([]byte(b), &ret)
				if err != nil {
					return errors.System.AddDetail(err)
				}
				rets = append(rets, ret)
			}
			m.Lock()
			defer m.Unlock()
			resp.Lists = append(resp.Lists, rets)
			return nil
		})
	}
	err = egg.Wait()
	return
}
