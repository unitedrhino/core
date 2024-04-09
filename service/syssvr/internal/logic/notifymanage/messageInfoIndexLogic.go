package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageInfoIndexLogic {
	return &MessageInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MessageInfoIndexLogic) MessageInfoIndex(in *sys.MessageInfoIndexReq) (*sys.MessageInfoIndexResp, error) {
	f := relationDB.MessageInfoFilter{Group: in.Group, NotifyCode: in.NotifyCode, IsDirectNotify: def.False}
	total, err := relationDB.NewMessageInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := relationDB.NewMessageInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	var list []*sys.MessageInfo
	for _, v := range pos {
		list = append(list, utils.Copy[sys.MessageInfo](v))
	}
	return &sys.MessageInfoIndexResp{Total: total, List: list}, nil
}
