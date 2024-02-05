package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/dict"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.WithID) (resp *types.DictInfo, err error) {
	ret, err := l.svcCtx.DictM.DictInfoRead(l.ctx, &sys.WithID{Id: req.ID})
	return dict.ToInfoTypes(ret), err
}
