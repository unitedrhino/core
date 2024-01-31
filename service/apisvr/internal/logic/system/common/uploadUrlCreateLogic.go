package common

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/oss"
	"gitee.com/i-Things/core/shared/oss/common"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadUrlCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadUrlCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadUrlCreateLogic {
	return &UploadUrlCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadUrlCreateLogic) UploadUrlCreate(req *types.UploadUrlCreateReq) (resp *types.UploadUrlCreateResp, err error) {

	filePath, err := oss.GetFilePath(&oss.SceneInfo{
		Business: req.Business,
		Scene:    req.Scene,
		FilePath: req.FilePath}, req.Rename)
	if err != nil {
		l.Errorf("%s.GetFilePath err:%v", utils.FuncName(), err)
		return nil, err
	}

	url, err := l.svcCtx.OssClient.TemporaryBucket().SignedPutUrl(l.ctx, filePath, int64(24*3600), common.OptionKv{})
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	resp = &types.UploadUrlCreateResp{
		FilePath:  filePath,
		UploadUri: url,
	}
	return resp, err
}
