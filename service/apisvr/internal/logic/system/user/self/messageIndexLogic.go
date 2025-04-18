package self

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageIndexLogic {
	return &MessageIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageIndexLogic) MessageIndex(req *types.UserMessageIndexReq) (resp *types.UserMessageIndexResp, err error) {
	ret, err := l.svcCtx.UserRpc.UserMessageIndex(l.ctx, utils.Copy[sys.UserMessageIndexReq](req))
	if err != nil {
		return nil, err
	}
	var list []*types.UserMessage
	for _, v := range ret.List {
		val := utils.Copy[types.UserMessage](v)
		val.MessageInfo = utils.Copy[types.MessageInfo](v.Message)
		list = append(list, val)
	}
	return &types.UserMessageIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     list,
	}, err
}
