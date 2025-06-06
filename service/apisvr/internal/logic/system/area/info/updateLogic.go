package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.AreaInfo) error {
	dmReq := &sys.AreaInfo{
		AreaID:             req.AreaID,
		ParentAreaID:       req.ParentAreaID,
		ProjectID:          req.ProjectID,
		AreaName:           req.AreaName,
		Position:           logic.ToSysPointRpc(req.Position),
		Desc:               utils.ToRpcNullString(req.Desc),
		UseBy:              req.UseBy,
		AreaImg:            req.AreaImg,
		IsUpdateAreaImg:    req.IsUpdateAreaImg,
		Tags:               req.Tags,
		ConfigFile:         req.ConfigFile,
		IsUpdateConfigFile: req.IsUpdateConfigFile,
	}
	_, err := l.svcCtx.AreaM.AreaInfoUpdate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AreaManage req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
