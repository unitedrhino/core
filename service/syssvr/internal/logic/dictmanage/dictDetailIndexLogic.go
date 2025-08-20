package dictmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailIndexLogic {
	return &DictDetailIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDetailIndexLogic) DictDetailIndex(in *sys.DictDetailIndexReq) (*sys.DictDetailIndexResp, error) {
	f := relationDB.DictDetailFilter{
		DictCode: in.DictCode,
		ParentID: in.ParentID,
		Status:   in.Status,
		Label:    in.Label,
		Value:    in.Value,
		Values:   in.Values,
	}
	repo := relationDB.NewDictDetailRepo(l.ctx)
	total, err := repo.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := repo.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultSort())
	if err != nil {
		return nil, err
	}
	return &sys.DictDetailIndexResp{Total: total, List: ToDictDetailsPb(pos)}, nil
}
