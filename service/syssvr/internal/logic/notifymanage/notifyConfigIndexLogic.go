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

type NotifyConfigIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigIndexLogic {
	return &NotifyConfigIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyConfigIndexLogic) NotifyConfigIndex(in *sys.NotifyConfigIndexReq) (*sys.NotifyConfigIndexResp, error) {
	db := relationDB.NewNotifyConfigRepo(l.ctx)
	f := relationDB.NotifyConfigFilter{
		Code:          in.Code,
		Group:         in.Group,
		Name:          in.Name,
		WithTemplates: in.WithChooseTemplateID,
	}
	totaol, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	var list []*sys.NotifyConfig
	for _, po := range pos {
		v := utils.Copy[sys.NotifyConfig](po)
		for _, val := range po.Templates {
			v.TemplateIDs = append(v.TemplateIDs, val.TemplateID)
		}
	}
	return &sys.NotifyConfigIndexResp{List: list, Total: totaol}, nil
}
