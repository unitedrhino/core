package areamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaProfileIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreaProfileIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaProfileIndexLogic {
	return &AreaProfileIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AreaProfileIndexLogic) AreaProfileIndex(in *sys.AreaProfileIndexReq) (*sys.AreaProfileIndexResp, error) {
	ret, err := relationDB.NewAreaProfileRepo(l.ctx).FindByFilter(l.ctx, relationDB.AreaProfileFilter{
		Codes:  in.Codes,
		AreaID: in.AreaID,
	}, nil)
	return &sys.AreaProfileIndexResp{Profiles: utils.CopySlice[sys.AreaProfile](ret)}, err

}
