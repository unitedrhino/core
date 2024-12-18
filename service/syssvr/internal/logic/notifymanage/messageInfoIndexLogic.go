package notifymanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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
	f := relationDB.MessageInfoFilter{Group: in.Group, NotifyCode: in.NotifyCode, IsDirectNotify: def.False, WithNotifyConfig: true}
	total, err := relationDB.NewMessageInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := relationDB.NewMessageInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "notifyTime",
		Sort:  stores.OrderDesc,
	}))
	if err != nil {
		return nil, err
	}
	var list []*sys.MessageInfo
	for _, v := range pos {
		do := utils.Copy[sys.MessageInfo](v)
		if v.NotifyConfig != nil {
			do.NotifyName = v.NotifyConfig.Name
		}
		list = append(list, do)
	}
	return &sys.MessageInfoIndexResp{Total: total, List: list}, nil
}
