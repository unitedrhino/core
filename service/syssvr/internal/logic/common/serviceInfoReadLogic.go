package commonlogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewServiceInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServiceInfoReadLogic {
	return &ServiceInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ServiceInfoReadLogic) ServiceInfoRead(in *sys.WithCode) (*sys.ServiceInfo, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ret, err := relationDB.NewServiceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ServiceInfoFilter{
		Code: in.Code,
	})
	if errors.Cmp(err, errors.NotFind) {
		return &sys.ServiceInfo{
			Code: in.Code,
		}, nil
	}
	return utils.Copy[sys.ServiceInfo](ret), err
}
