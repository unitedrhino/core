package notifymanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyTemplateIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyTemplateIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyTemplateIndexLogic {
	return &NotifyTemplateIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyTemplateIndexLogic) NotifyTemplateIndex(in *sys.NotifyTemplateIndexReq) (*sys.NotifyTemplateIndexResp, error) {
	db := relationDB.NewNotifyTemplateRepo(l.ctx)
	f := relationDB.NotifyTemplateFilter{
		Name:       in.Name,
		NotifyCode: in.NotifyCode,
		Type:       in.Type,
	}
	total, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	return &sys.NotifyTemplateIndexResp{Total: total, List: utils.CopySlice[sys.NotifyTemplate](pos)}, nil
}
