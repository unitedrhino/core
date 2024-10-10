package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaApplyIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaApplyIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaApplyIndexLogic {
	return &UserAreaApplyIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaApplyIndexLogic) UserAreaApplyIndex(in *sys.UserAreaApplyIndexReq) (*sys.UserAreaApplyIndexResp, error) {
	f := relationDB.UserAreaApplyFilter{AuthTypes: in.AuthTypes, AreaID: in.AreaID}
	total, err := relationDB.NewUserAreaApplyRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	list, err := relationDB.NewUserAreaApplyRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	return &sys.UserAreaApplyIndexResp{List: ToUserAreaApplyInfos(list), Total: total}, nil
}
