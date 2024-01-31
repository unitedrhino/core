package accessmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessInfoUpdateLogic {
	return &AccessInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccessInfoUpdateLogic) AccessInfoUpdate(in *sys.AccessInfo) (*sys.Response, error) {
	old, err := relationDB.NewAccessRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.Name = in.Name
	old.Group = in.Group
	old.IsNeedAuth = in.IsNeedAuth
	old.Desc = in.Desc
	err = relationDB.NewAccessRepo(l.ctx).Update(l.ctx, old)
	return &sys.Response{}, err
}
