package common

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/result"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadFileLogic {
	return &DownloadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadFileLogic) DownloadFile(req *types.DownloadFileReq, w http.ResponseWriter, r *http.Request) {
	ossLocal, ok := l.svcCtx.OssClient.Handle.(*oss.Local)
	if !ok {
		result.Http(w, r, nil, errors.System.AddMsgf("本地存储不支持"))
		return
	}
	err := ossLocal.DownloadFile(l.ctx, req.FilePath, req.Sign, w)
	if err != nil {
		result.Http(w, r, nil, err)
		return
	}

	return
}
