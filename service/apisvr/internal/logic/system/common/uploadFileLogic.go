package common

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/shared/oss"
	"gitee.com/i-Things/core/shared/oss/common"
	"gitee.com/i-Things/core/shared/utils"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadFileLogic) UploadFile() (resp *types.UploadFileResp, err error) {
	file, fh, err := l.r.FormFile("file")
	if err != nil {
		return resp, err
	}
	defer file.Close()
	newFilePath, err := oss.GetFilePath2(l.ctx, fh)
	if err != nil {
		l.Errorf("%s.GetFilePath err:%v", utils.FuncName(), err)
		return nil, err
	}
	fileUri, err := l.svcCtx.OssClient.TemporaryBucket().Upload(l.ctx, newFilePath, file, common.OptionKv{})
	if err != nil {
		return resp, err
	}
	return &types.UploadFileResp{
		FileUri:  fileUri,
		FilePath: newFilePath,
	}, err
}
