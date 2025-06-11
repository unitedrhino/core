package areamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaProfileReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreaProfileReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaProfileReadLogic {
	return &AreaProfileReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AreaProfileReadLogic) AreaProfileRead(in *sys.AreaProfileReadReq) (*sys.AreaProfile, error) {
	ret, err := relationDB.NewAreaProfileRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.AreaProfileFilter{
		Code:   in.Code,
		AreaID: in.AreaID,
	})
	if errors.Cmp(err, errors.NotFind) {
		return &sys.AreaProfile{
			AreaID: in.AreaID,
			Code:   in.Code,
			Params: "",
		}, nil
	}
	return utils.Copy[sys.AreaProfile](ret), err
}
