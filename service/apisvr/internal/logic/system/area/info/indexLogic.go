package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.AreaInfoIndexReq) (resp *types.AreaInfoIndexResp, err error) {
	dmReq := &sys.AreaInfoIndexReq{
		Page:         logic.ToSysPageRpc(req.Page),
		ProjectID:    req.ProjectID,
		AreaIDs:      req.AreaIDs,
		ParentAreaID: req.ParentAreaID,
	}
	dmResp, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AreaManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	return &types.AreaInfoIndexResp{
		Total: dmResp.Total,
		List:  ToAreaInfosTypes(dmResp.List),
	}, nil
}
