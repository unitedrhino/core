package dictmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictInfoUpdateLogic {
	return &DictInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictInfoUpdateLogic) DictInfoUpdate(in *sys.DictInfo) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	repo := relationDB.NewDictInfoRepo(l.ctx)
	old, err := repo.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	old.Name = in.Name
	old.Group = in.Group
	old.Desc = in.Desc.GetValue()
	old.Body = in.Body.GetValue()
	err = repo.Update(l.ctx, old)
	return &sys.Empty{}, err
}
