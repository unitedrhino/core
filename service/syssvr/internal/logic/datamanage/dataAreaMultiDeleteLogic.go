package datamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataAreaMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataAreaMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataAreaMultiDeleteLogic {
	return &DataAreaMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DataAreaMultiDeleteLogic) DataAreaMultiDelete(in *sys.DataAreaMultiDeleteReq) (*sys.Empty, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID != 0 {
		uc.ProjectID = in.ProjectID
	} else {
		in.ProjectID = uc.ProjectID
	}
	err := relationDB.NewDataAreaRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DataAreaFilter{Targets: []*relationDB.Target{{Type: in.TargetType, ID: in.TargetID}}, AreaIDs: in.AreaIDs})
	return &sys.Empty{}, err
}
