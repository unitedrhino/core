package commonlogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewServiceInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServiceInfoUpdateLogic {
	return &ServiceInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ServiceInfoUpdateLogic) ServiceInfoUpdate(in *sys.ServiceInfo) (*sys.Empty, error) {
	old, err := relationDB.NewServiceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ServiceInfoFilter{
		Code: in.Code,
	})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		old = &relationDB.SysServiceInfo{Code: in.Code}
	}
	old.Name = in.Name
	old.Desc = in.Desc
	old.Version = in.Version
	err = relationDB.NewServiceInfoRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
