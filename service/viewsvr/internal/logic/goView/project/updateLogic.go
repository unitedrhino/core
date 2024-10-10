package project

import (
	"context"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/domain"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/oss/common"

	"gitee.com/unitedrhino/core/service/viewsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ProjectInfo) error {
	pi, err := relationDB.NewProjectInfoRepo(l.ctx).FindOne(l.ctx, req.ID)
	if err != nil {
		return err
	}
	if req.Name != "" {
		pi.Name = req.Name
	}
	if req.Desc != "" {
		pi.Desc = req.Desc
	}
	if req.IndexImage != "" && req.IndexImage != pi.IndexImage {
		if pi.IndexImage != "" {
			err := l.svcCtx.OssClient.PublicBucket().Delete(l.ctx, pi.IndexImage, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", pi.IndexImage, err)
			}
		}
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, domain.BusinessView, domain.SceneProjectIndexImage, oss.GetFileNameWithPath(req.IndexImage))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(req.IndexImage, nwePath)
		if err != nil {
			return errors.System.AddDetail(err)
		}
		pi.IndexImage = path
	}
	if req.Status != 0 {
		pi.Status = req.Status
	}
	err = relationDB.NewProjectInfoRepo(l.ctx).Update(l.ctx, pi)
	return err
}
