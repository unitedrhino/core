package projectmanagelogic

import (
	"context"
	"fmt"
	"strings"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/oss"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectCrudUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProjectCrudUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectCrudUpdateLogic {
	return &ProjectCrudUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新项目crud
func (l *ProjectCrudUpdateLogic) ProjectCrudUpdate(in *sys.ProjectCrud) (*sys.Empty, error) {
	old, err := relationDB.NewProjectCurlRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Params != nil {
		var param = map[string]string{}
		for k, v := range in.Params {
			if !(strings.HasSuffix(k, "Img") || strings.HasSuffix(k, "File")) {
				param[k] = v
				continue
			}
			if !l.svcCtx.OssClient.IsFilePath(v) {
				param[k] = v
				continue
			}
			nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProject, "curl",
				fmt.Sprintf("%d/%s/%d/%s", old.ProjectID, old.Purpose, old.ID, oss.GetFileNameWithPath(v)))
			path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(v, nwePath)
			if err != nil {
				l.Error(err)
			} else {
				param[k] = path
			}
		}
		old.Params = param
	}
	if in.Sort != 0 {
		old.Sort = in.Sort
	}
	err = relationDB.NewProjectCurlRepo(l.ctx).Update(l.ctx, old)

	return &sys.Empty{}, err
}
