package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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

func (l *DictInfoUpdateLogic) DictInfoUpdate(in *sys.DictInfo) (*sys.Response, error) {
	repo := relationDB.NewDictInfoRepo(l.ctx)
	old, err := repo.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.Name = in.Name
	old.Type = in.Type
	old.Status = in.Status
	old.Desc = in.Desc.GetValue()
	old.Body = in.Body.GetValue()
	err = repo.Update(l.ctx, old)
	return &sys.Response{}, err
}
