package common

import (
	"context"
	"net/http"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DebugGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDebugGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DebugGetLogic {
	return &DebugGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DebugGetLogic) DebugGet(r *http.Request) (resp *types.DebugResp, err error) {
	var headers = make(map[string]string)
	for k, v := range r.Header {
		headers[k] = v[0]
	}
	return &types.DebugResp{
		RequestUri: r.RequestURI,
		Headers:    headers,
	}, nil
}
