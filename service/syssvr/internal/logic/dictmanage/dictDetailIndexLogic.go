package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

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
	f := relationDB.DictDetailFilter{DictID: in.DictID}
	repo := relationDB.NewDictDetailRepo(l.ctx)
	total, err := repo.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := repo.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	var list []*sys.DictDetail
	for _, v := range pos {
		list = append(list, &sys.DictDetail{
			Id:     v.ID,
			DictID: v.DictID,
			Label:  v.Label,
			Value:  v.Value,
			Extend: v.Extend,
			Sort:   v.Sort,
			Desc:   utils.ToRpcNullString(v.Desc),
			Status: v.Status,
			Body:   utils.ToRpcNullString(v.Body),
		})
	}
	return &sys.DictDetailIndexResp{Total: total, List: ToDictDetailsPb(pos)}, nil
}
