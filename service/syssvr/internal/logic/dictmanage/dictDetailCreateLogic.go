package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailCreateLogic {
	return &DictDetailCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDetailCreateLogic) DictDetailCreate(in *sys.DictDetail) (*sys.WithID, error) {
	_, err := relationDB.NewApiInfoRepo(l.ctx).FindOne(l.ctx, in.DictID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("字典未定义")
		}
		return nil, err
	}
	po := relationDB.SysDictDetail{
		DictID: in.DictID,
		Label:  in.Label,
		Value:  in.Value,
		Extend: in.Extend,
		Status: in.Status,
		Sort:   in.Sort,
		Desc:   in.Desc.GetValue(),
		Body:   in.Body.GetValue(),
	}
	err = relationDB.NewDictDetailRepo(l.ctx).Insert(l.ctx, &po)
	return &sys.WithID{Id: po.ID}, err
}
