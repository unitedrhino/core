package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/dict"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

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

func (l *CreateLogic) Create(req *types.DictInfo) (resp *types.WithID, err error) {
	ret, err := l.svcCtx.DictM.DictInfoCreate(l.ctx, dict.ToInfoPb(req))
	if err != nil {
		return nil, err
	}
	return &types.WithID{ID: ret.Id}, nil
}
