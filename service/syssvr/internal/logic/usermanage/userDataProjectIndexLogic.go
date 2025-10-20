package usermanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	datamanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/datamanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDataProjectIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UapDB *relationDB.DataProjectRepo
}

func NewUserDataProjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDataProjectIndexLogic {
	return &UserDataProjectIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UapDB:  relationDB.NewDataProjectRepo(ctx),
	}
}

func (l *UserDataProjectIndexLogic) UserDataProjectIndex(in *sys.UserDataProjectIndexReq) (*sys.UserDataProjectIndexResp, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}

	filter := relationDB.DataProjectFilter{
		Targets: []*relationDB.Target{{Type: def.TargetUser, ID: in.UserID}},
	}
	rs, err := relationDB.NewUserRoleRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserRoleFilter{UserID: in.UserID}, nil)
	if err != nil {
		return nil, err
	}
	for _, v := range rs {
		filter.Targets = append(filter.Targets, &relationDB.Target{Type: def.TargetUser, ID: v.ID})
	}
	total, err := l.UapDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	poArr, err := l.UapDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}

	list := make([]*sys.DataProject, 0, len(poArr))
	for _, po := range poArr {
		list = append(list, datamanagelogic.ProjectPoToPb(po))
	}
	return &sys.UserDataProjectIndexResp{List: list, Total: total}, nil
}
