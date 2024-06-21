package areamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
			Code:   in.Code,
			Params: "",
		}, nil
	}
	return utils.Copy[sys.AreaProfile](ret), err
}
