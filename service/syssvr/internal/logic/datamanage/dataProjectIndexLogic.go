package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataProjectIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UapDB *relationDB.DataProjectRepo
}

func NewDataProjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataProjectIndexLogic {
	return &DataProjectIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UapDB:  relationDB.NewDataProjectRepo(ctx),
	}
}

func (l *DataProjectIndexLogic) DataProjectIndex(in *sys.DataProjectIndexReq) (*sys.DataProjectIndexResp, error) {
	var (
		list  []*sys.DataProject
		total int64
		err   error
	)
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID != 0 {
		uc.ProjectID = in.ProjectID
	} else {
		in.ProjectID = uc.ProjectID
	}
	if !uc.IsAdmin && in.TargetType != def.TargetUser {
		return nil, errors.Permissions.AddMsg("非管理员只能获取用户类型的")
	}
	if uc.AllProject {
		in.ProjectID = 0
	}
	if in.ProjectID != 0 {
		if uc.IsAdmin || uc.ProjectAuth[in.ProjectID] != nil {
			in.ProjectID = in.ProjectID
		}
	}
	filter := relationDB.DataProjectFilter{
		ProjectID: in.ProjectID,
		AuthType:  in.AuthType,
		Target: &relationDB.Target{
			Type: in.TargetType,
			ID:   in.TargetID,
		},
	}
	if in.ProjectID != 0 && uc.IsAdmin {
		filter.ProjectID = in.ProjectID
	}
	total, err = l.UapDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	poArr, err := l.UapDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}

	list = make([]*sys.DataProject, 0, len(poArr))
	for _, po := range poArr {
		list = append(list, transProjectPoToPb(po))
	}
	return &sys.DataProjectIndexResp{List: list, Total: total}, nil
}
