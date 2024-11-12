package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.ProjectWithID) (resp *types.ProjectInfo, err error) {
	dmResp, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: req.ProjectID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ProjectManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	var user *sys.UserInfo
	if req.WithAdminUser {
		user, err = l.svcCtx.UserCache.GetData(l.ctx, dmResp.AdminUserID)
		if err != nil {
			l.Error(err)
		}
	}

	return system.ProjectInfoToApi(dmResp, user), nil
}
