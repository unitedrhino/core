package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyChannelIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyChannelIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyChannelIndexLogic {
	return &NotifyChannelIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyChannelIndexLogic) NotifyChannelIndex(in *sys.NotifyChannelIndexReq) (*sys.NotifyChannelIndexResp, error) {
	db := relationDB.NewNotifyChannelRepo(l.ctx)
	f := relationDB.NotifyChannelFilter{
		Name: in.Name,
		Type: in.Type,
	}
	total, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	return &sys.NotifyChannelIndexResp{Total: total, List: utils.CopySlice[sys.NotifyChannel](pos)}, nil

}
