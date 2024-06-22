package profile

import (
	"context"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

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

func (l *ProfileIndexLogic) ProfileIndex(req *types.ProjectProfileIndexReq) (resp *types.ProjectProfileIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}
