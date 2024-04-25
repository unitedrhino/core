package project

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectIndexLogic {
	return &ProjectIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectIndexLogic) ProjectIndex(req *types.DataProjectIndexReq) (resp *types.DataProjectIndexResp, err error) {
	dto := &sys.DataProjectIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		ProjectID:  req.ProjectID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
	}
	dmResp, err := l.svcCtx.DataM.DataProjectIndex(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.DataProjectIndex req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	list := ToProjectApis(dmResp.List)
	return &types.DataProjectIndexResp{
		Total: dmResp.Total,
		List:  list,
	}, nil
}
