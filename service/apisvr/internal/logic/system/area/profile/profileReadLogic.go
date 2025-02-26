package profile

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileReadLogic {
	return &ProfileReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileReadLogic) ProfileRead(req *types.AreaProfileReadReq) (resp *types.AreaProfile, err error) {
	ret, err := l.svcCtx.AreaM.AreaProfileRead(l.ctx, utils.Copy[sys.AreaProfileReadReq](req))
	return utils.Copy[types.AreaProfile](ret), err
}
