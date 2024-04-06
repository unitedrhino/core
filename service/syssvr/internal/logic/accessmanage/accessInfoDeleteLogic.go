package accessmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessInfoDeleteLogic {
	return &AccessInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccessInfoDeleteLogic) AccessInfoDelete(in *sys.WithID) (*sys.Empty, error) {
	err := relationDB.NewAccessRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Empty{}, err
}
