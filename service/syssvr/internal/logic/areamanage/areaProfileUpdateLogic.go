package areamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaProfileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreaProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaProfileUpdateLogic {
	return &AreaProfileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AreaProfileUpdateLogic) AreaProfileUpdate(in *sys.AreaProfile) (*sys.Empty, error) {
	old, err := relationDB.NewAreaProfileRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.AreaProfileFilter{
		Code:   in.Code,
		AreaID: in.AreaID,
	})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		old = &relationDB.SysAreaProfile{AreaID: stores.AreaID(in.AreaID), Code: in.Code}
	}
	old.Params = in.Params
	err = relationDB.NewAreaProfileRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}