package profile

import (
	"context"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileUpdateLogic {
	return &ProfileUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileUpdateLogic) ProfileUpdate(req *types.ProjectProfile) error {
	// todo: add your logic here and delete this line

	return nil
}
