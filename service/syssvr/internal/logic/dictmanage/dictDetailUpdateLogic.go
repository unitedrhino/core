package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailUpdateLogic {
	return &DictDetailUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDetailUpdateLogic) DictDetailUpdate(in *sys.DictDetail) (*sys.Response, error) {
	old, err := relationDB.NewDictDetailRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.Label = in.Label
	old.Value = in.Value
	old.Extend = in.Extend
	old.Status = in.Status
	old.Sort = in.Sort
	old.Desc = in.Desc.GetValue()
	old.Body = in.Body.GetValue()
	err = relationDB.NewDictDetailRepo(l.ctx).Update(l.ctx, old)
	return &sys.Response{}, err
}