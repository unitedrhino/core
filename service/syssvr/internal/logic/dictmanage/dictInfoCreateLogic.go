package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictInfoCreateLogic {
	return &DictInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictInfoCreateLogic) DictInfoCreate(in *sys.DictInfo) (*sys.WithID, error) {
	po := relationDB.SysDictInfo{
		Name:   in.Name,
		Type:   in.Type,
		Status: in.Status,
		Desc:   in.Desc.GetValue(),
		Body:   in.Body.GetValue(),
	}
	err := relationDB.NewDictInfoRepo(l.ctx).Insert(l.ctx, &po)
	return &sys.WithID{Id: po.ID}, err
}
