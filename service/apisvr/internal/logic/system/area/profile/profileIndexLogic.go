package profile

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileIndexLogic {
	return &ProfileIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileIndexLogic) ProfileIndex(req *types.AreaProfileIndexReq) (resp *types.AreaProfileIndexResp, err error) {
	ret, err := l.svcCtx.AreaM.AreaProfileIndex(l.ctx, utils.Copy[sys.AreaProfileIndexReq](req))

	return utils.Copy[types.AreaProfileIndexResp](ret), nil
}
