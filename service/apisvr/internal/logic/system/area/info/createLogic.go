package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.AreaInfo) (*types.AreaWithID, error) {
	if req.AreaName == "" || ////root节点不为0
		req.ParentAreaID == def.NotClassified { //未分类不能有下属的区域
		return nil, errors.Parameter
	}
	dmReq := &sys.AreaInfo{
		ParentAreaID:    req.ParentAreaID,
		ProjectID:       req.ProjectID,
		AreaName:        req.AreaName,
		Position:        logic.ToSysPointRpc(req.Position),
		Desc:            utils.ToRpcNullString(req.Desc),
		UseBy:           req.UseBy,
		AreaImg:         req.AreaImg,
		IsUpdateAreaImg: req.IsUpdateAreaImg,
		IsSysCreated:    req.IsSysCreated,
	}
	resp, err := l.svcCtx.AreaM.AreaInfoCreate(l.ctx, dmReq)
	if er := errors.Fmt(err); er != nil {
		l.Errorf("%s.rpc.AreaManage req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.AreaWithID{AreaID: resp.AreaID}, nil
}
