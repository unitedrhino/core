package common

import (
	"context"
	"io"
	"net/http"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DebugPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDebugPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DebugPostLogic {
	return &DebugPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DebugPostLogic) DebugPost(r *http.Request) (resp *types.DebugResp, err error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logx.Error(err)
	}
	var headers = make(map[string]string)
	for k, v := range r.Header {
		headers[k] = v[0]
	}
	return &types.DebugResp{
		RequestUri: r.RequestURI,
		Headers:    headers,
		Body:       string(body),
	}, nil
}
