package projectmanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"

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
	if in.Params != "" {
		var oldParams = make(map[string]interface{})
		err = json.Unmarshal([]byte(in.Params), &oldParams)
		if err != nil {
			return nil, errors.Parameter.AddMsg("params 不是json").AddDetail(err)
		}
		var param = map[string]any{}
		for k, v := range oldParams {
			if !(strings.HasSuffix(k, "Img") || strings.HasSuffix(k, "File")) {
				param[k] = v
				continue
			}
			if !l.svcCtx.OssClient.IsFilePath(cast.ToString(v)) {
				param[k] = v
				continue
			}
			nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProject, "curl",
				fmt.Sprintf("%d/%s/%d/%s", old.ProjectID, old.Purpose, old.ID, oss.GetFileNameWithPath(cast.ToString(v))))
			path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(cast.ToString(v), nwePath)
			if err != nil {
				l.Error(err)
			} else {
				param[k] = path
			}
		}
		old.Params = utils.MarshalNoErr(param)
	}
	if in.Sort != 0 {
		old.Sort = in.Sort
	}
	err = relationDB.NewProjectCurlRepo(l.ctx).Update(l.ctx, old)

	return &sys.Empty{}, err
}
