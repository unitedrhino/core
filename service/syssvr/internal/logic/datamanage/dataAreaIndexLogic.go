package datamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataAreaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UaaDB *relationDB.DataAreaRepo
}

func NewDataAreaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataAreaIndexLogic {
	return &DataAreaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UaaDB:  relationDB.NewDataAreaRepo(ctx),
	}
}

func (l *DataAreaIndexLogic) DataAreaIndex(in *sys.DataAreaIndexReq) (*sys.DataAreaIndexResp, error) {
	var (
		list  []*sys.DataArea
		total int64
		err   error
	)
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID != 0 {
		uc.ProjectID = in.ProjectID
	} else {
		in.ProjectID = uc.ProjectID
	}
	filter := relationDB.DataAreaFilter{
		Targets:   []*relationDB.Target{{Type: in.TargetType, ID: in.TargetID}},
		ProjectID: in.ProjectID,
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
		list = append(list, transAreaPoToPb(po))
	}
	return &sys.DataAreaIndexResp{List: list, Total: total}, nil
}
