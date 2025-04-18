package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.ProjectInfoIndexReq) (resp *types.ProjectInfoIndexResp, err error) {
	dmReq := &sys.ProjectInfoIndexReq{
		Page:         logic.ToSysPageRpc(req.Page),
		ProjectName:  req.ProjectName,
		ProjectIDs:   req.ProjectIDs,
		IsGetAll:     req.IsGetAll,
		WithTopAreas: req.WithTopAreas,
		TenantCode:   req.TenantCode,
	}
	dmResp, err := l.svcCtx.ProjectM.ProjectInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProjectManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	list := make([]*types.ProjectInfo, 0, len(dmResp.List))
	for _, pb := range dmResp.List {
		var user *sys.UserInfo
		if req.WithAdminUser {
			user, err = l.svcCtx.UserCache.GetData(l.ctx, pb.AdminUserID)
			if err != nil {
				l.Error(err)
			}
		}
		list = append(list, system.ProjectInfoToApi(pb, user))
	}

	return &types.ProjectInfoIndexResp{
		PageResp: logic.ToPageResp(req.Page, dmResp.Total),
		List:     list,
	}, nil
}
