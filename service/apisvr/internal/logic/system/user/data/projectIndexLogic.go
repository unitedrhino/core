package data

import (
	"context"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取项目权限列表
func NewProjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectIndexLogic {
	return &ProjectIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectIndexLogic) ProjectIndex(req *types.UserDataProjectIndexReq) (resp *types.UserDataProjectIndexResp, err error) {
	ret, err := l.svcCtx.UserRpc.UserDataProjectIndex(l.ctx, utils.Copy[sys.UserDataProjectIndexReq](req))
	if err != nil {
		l.Errorf("%s.rpc.DataProjectIndex req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	svcCtx := l.svcCtx

	list := ToProjectApis(l.ctx, svcCtx, ret.List)
	return &types.UserDataProjectIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     list,
	}, nil
}

func ToProjectApis(ctx context.Context, svcCtx *svc.ServiceContext, in []*sys.DataProject) (ret []*types.UserDataProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ui := &types.ProjectInfo{}
		if svcCtx != nil {
			u, err := svcCtx.ProjectCache.GetData(ctx, v.ProjectID)
			if err != nil {
				continue
			}
			ui = utils.Copy[types.ProjectInfo](u)
		}
		ret = append(ret, &types.UserDataProject{ProjectID: v.ProjectID, AuthType: v.AuthType, TargetID: v.TargetID, UpdatedTime: v.UpdatedTime, Project: ui})
	}
	return
}
