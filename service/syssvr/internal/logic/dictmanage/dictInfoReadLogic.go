package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

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

func (l *DictInfoReadLogic) DictInfoRead(in *sys.WithID) (*sys.DictInfo, error) {
	one, err := relationDB.NewDictInfoRepo(l.ctx).FindOneByFilter(l.ctx,
		relationDB.DictInfoFilter{ID: in.Id, WithDetails: true})
	if err != nil {
		return nil, err
	}
	return ToDictInfoPb(one), nil
}
