package notifymanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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
		Code:  in.Code,
		Group: in.Group,
		Name:  in.Name,
	}
	totaol, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "createdTime",
		Sort:  stores.OrderDesc,
	}))
	if err != nil {
		return nil, err
	}
	var list []*sys.NotifyConfig
	for _, po := range pos {
		v := utils.Copy[sys.NotifyConfig](po)
		list = append(list, v)
	}
	return &sys.NotifyConfigIndexResp{List: list, Total: totaol}, nil
}
