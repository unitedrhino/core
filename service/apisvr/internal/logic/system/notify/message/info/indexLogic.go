package info

import (
	"context"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.MessageInfoIndexReq) (resp *types.MessageInfoIndexResp, err error) {
	ret, err := l.svcCtx.NotifyM.MessageInfoIndex(l.ctx, utils.Copy[sys.MessageInfoIndexReq](req))
	if err != nil {
		return nil, err
	}
	return &types.MessageInfoIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     utils.CopySlice[types.MessageInfo](ret.List),
	}, nil
}
