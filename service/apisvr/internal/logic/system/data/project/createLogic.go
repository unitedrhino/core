package project

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.DataProjectSaveReq) (resp *types.DataProject, err error) {
	_, err = l.svcCtx.DataM.DataProjectCreate(l.ctx, utils.Copy[sys.DataProjectSaveReq](req))
	if err != nil {
		return nil, err
	}
	resp = utils.Copy[types.DataProject](req)
	if req.TargetType == def.TargetUser {
		u, err := l.svcCtx.UserCache.GetData(l.ctx, req.TargetID)
		if err != nil {
			return nil, err
		}
		resp.User = utils.Copy[types.UserCore](u)
	}
	return
}
