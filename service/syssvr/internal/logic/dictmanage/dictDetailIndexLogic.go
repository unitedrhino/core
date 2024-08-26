package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
	}
	repo := relationDB.NewDictDetailRepo(l.ctx)
	total, err := repo.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := repo.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	return &sys.DictDetailIndexResp{Total: total, List: ToDictDetailsPb(pos)}, nil
}
