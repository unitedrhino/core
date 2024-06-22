package info

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

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

func (l *ProfileReadLogic) ProfileRead(req *types.WithCode) (resp *types.ProjectProfile, err error) {
	ret, err := l.svcCtx.ProjectM.ProjectProfileRead(l.ctx, utils.Copy[sys.ProjectProfileReadReq](req))
	return utils.Copy[types.ProjectProfile](ret), err
}
