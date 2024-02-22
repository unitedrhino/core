package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/def"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictInfoReadLogic {
	return &DictInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictInfoReadLogic) DictInfoRead(in *sys.DictInfoReadReq) (*sys.DictInfo, error) {
	var (
		po  *relationDB.SysDictInfo
		err error
	)
	switch in.Id {
	case def.RootNode, 0:
		po = &relationDB.SysDictInfo{
			ID:   def.RootNode,
			Name: "全部字典",
		}
	default:
		po, err = relationDB.NewDictInfoRepo(l.ctx).FindOneByFilter(l.ctx,
			relationDB.DictInfoFilter{ID: in.Id, WithDetails: in.WithDetails})
		if err != nil {
			return nil, err
		}
	}
	if !in.WithChildren {
		return ToDictInfoPb(po, nil), nil
	}
	children, err := relationDB.NewDictInfoRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.DictInfoFilter{IDPath: po.IDPath, WithDetails: in.WithDetails}, nil)
	if err != nil {
		return nil, err
	}

	return ToDictInfoPb(po, children), nil
}
