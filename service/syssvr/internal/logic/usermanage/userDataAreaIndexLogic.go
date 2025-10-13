package usermanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	datamanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/datamanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDataAreaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UaaDB *relationDB.DataAreaRepo
}

func NewUserDataAreaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDataAreaIndexLogic {
	return &UserDataAreaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UaaDB:  relationDB.NewDataAreaRepo(ctx),
	}
}

func (l *UserDataAreaIndexLogic) UserDataAreaIndex(in *sys.UserDataAreaIndexReq) (*sys.UserDataAreaIndexResp, error) {
	var (
		list  []*sys.DataArea
		total int64
		err   error
	)
	filter := relationDB.DataAreaFilter{
		Targets:   []*relationDB.Target{{Type: def.TargetUser, ID: in.UserID}},
		ProjectID: in.ProjectID,
	}
	rs, err := relationDB.NewUserRoleRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserRoleFilter{UserID: in.UserID}, nil)
	if err != nil {
		return nil, err
	}
	for _, v := range rs {
		filter.Targets = append(filter.Targets, &relationDB.Target{Type: def.TargetUser, ID: v.ID})
	}

	total, err = l.UaaDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	poArr, err := l.UaaDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}

	list = make([]*sys.DataArea, 0, len(poArr))
	for _, po := range poArr {
		list = append(list, datamanagelogic.AreaPoToPb(po))
	}
	return &sys.UserDataAreaIndexResp{List: list, Total: total}, nil
}
