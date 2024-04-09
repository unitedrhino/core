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

type NotifyInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyInfoIndexLogic {
	return &NotifyInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyInfoIndexLogic) NotifyInfoIndex(in *sys.NotifyInfoIndexReq) (*sys.NotifyInfoIndexResp, error) {
	db := relationDB.NewNotifyInfoRepo(l.ctx)
	f := relationDB.NotifyInfoFilter{
		Code:  in.Code,
		Group: in.Group,
		Name:  in.Name,
	}
	totaol, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	list := utils.CopySlice[sys.NotifyInfo](pos)
	return &sys.NotifyInfoIndexResp{List: list, Total: totaol}, nil
}
