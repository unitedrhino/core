package dictmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictInfoIndexLogic {
	return &DictInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictInfoIndexLogic) DictInfoIndex(in *sys.DictInfoIndexReq) (*sys.DictInfoIndexResp, error) {
	f := relationDB.DictInfoFilter{
		Name:  in.Name,
		Group: in.Group,
	}
	total, err := relationDB.NewDictInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := relationDB.NewDictInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	var list []*sys.DictInfo
	for _, v := range pos {
		list = append(list, ToDictInfoPb(v))
	}
	return &sys.DictInfoIndexResp{Total: total, List: list}, nil
}
